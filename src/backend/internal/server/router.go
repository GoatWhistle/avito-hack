package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/shared/middleware"
)

const (
	readyPingTimeout = 2 * time.Second
	uploadPathParts  = 2
)

type ModuleRegistrar interface {
	RegisterRoutes(r chi.Router)
}

type uploadGuard func(ctx context.Context, itemID uuid.UUID) (bool, error)

type RouterDeps struct {
	Config    config.Config
	Pool      *pgxpool.Pool
	Logger    *slog.Logger
	Modules   []ModuleRegistrar
	WebSocket http.Handler
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(middleware.Recoverer(deps.Logger))
	r.Use(middleware.AccessLog(deps.Logger))
	r.Use(middleware.Metrics())
	r.Use(middleware.CORS(deps.Config.AllowedOrigins))

	r.Get("/healthz", healthHandler)
	r.Get("/readyz", readyHandler(deps.Pool))
	r.Handle("/metrics", promhttp.Handler())

	mountUploads(r, deps.Config.UploadURL, deps.Config.UploadDir, newUploadGuard(deps.Pool))

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Handle("/ws", deps.WebSocket)
		v1.Group(func(api chi.Router) {
			if deps.Config.RequestTimeout > 0 {
				api.Use(chimw.Timeout(deps.Config.RequestTimeout))
			}

			for _, module := range deps.Modules {
				module.RegisterRoutes(api)
			}
		})
	})

	return r
}

func mountUploads(r chi.Router, publicURL, dir string, guard uploadGuard) {
	if publicURL == "" || dir == "" {
		return
	}

	prefix := "/" + strings.Trim(publicURL, "/")
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))

	r.Get(prefix+"/*", func(w http.ResponseWriter, req *http.Request) {
		itemID, ok := uploadItemID(strings.TrimPrefix(req.URL.Path, prefix))
		if !ok {
			http.NotFound(w, req)

			return
		}

		visible, err := guard(req.Context(), itemID)
		if err != nil {
			writeStatus(w, http.StatusServiceUnavailable, `{"status":"unavailable"}`)

			return
		}

		if !visible {
			http.NotFound(w, req)

			return
		}

		w.Header().Set("Cache-Control", "private, max-age=86400")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fileServer.ServeHTTP(w, req)
	})
}

func uploadItemID(rest string) (uuid.UUID, bool) {
	parts := strings.Split(strings.Trim(path.Clean(rest), "/"), "/")
	if len(parts) != uploadPathParts {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

func newUploadGuard(pool *pgxpool.Pool) uploadGuard {
	if pool == nil {
		return func(context.Context, uuid.UUID) (bool, error) { return true, nil }
	}

	return func(ctx context.Context, itemID uuid.UUID) (bool, error) {
		const query = `
			SELECT EXISTS (
				SELECT 1 FROM items
				WHERE id = $1 AND deleted_at IS NULL AND status IN ('published', 'sold')
			)`

		var visible bool
		if err := pool.QueryRow(ctx, query, itemID).Scan(&visible); err != nil {
			return false, fmt.Errorf("check item visibility: %w", err)
		}

		return visible, nil
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeStatus(w, http.StatusOK, `{"status":"ok"}`)
}

func readyHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			writeStatus(w, http.StatusServiceUnavailable, `{"status":"unavailable"}`)

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), readyPingTimeout)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			writeStatus(w, http.StatusServiceUnavailable, `{"status":"unavailable"}`)

			return
		}

		writeStatus(w, http.StatusOK, `{"status":"ready"}`)
	}
}

func writeStatus(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(body)); err != nil {
		slog.Error("write response", slog.String("error", err.Error()))
	}
}
