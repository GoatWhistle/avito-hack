package server

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/middleware"
)

const (
	readyPingTimeout = 2 * time.Second
	uploadPathParts  = 2
	authPathPrefix   = "/api/v1/auth/"
	compressionLevel = 5
)

type ModuleRegistrar interface {
	RegisterRoutes(r chi.Router)
}

type uploadGuard func(ctx context.Context, displayID string) (uuid.UUID, bool, error)

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

	r.Group(func(metricsRoute chi.Router) {
		if deps.Config.RateLimitEnabled {
			metricsRoute.Use(middleware.RateLimit(middleware.RateLimitConfig{
				Rate:       deps.Config.RateLimitRPS,
				Burst:      deps.Config.RateLimitBurst,
				TrustProxy: deps.Config.RateLimitTrustProxy,
			}))
		}
		metricsRoute.Use(metricsAuth(deps.Config.MetricsToken))
		metricsRoute.Handle("/metrics", promhttp.Handler())
	})

	mountUploads(r, deps.Config.UploadURL, deps.Config.UploadDir, newUploadGuard(deps.Pool))

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Group(func(ws chi.Router) {
			if deps.Config.RateLimitEnabled {
				ws.Use(middleware.RateLimit(middleware.RateLimitConfig{
					Rate:       deps.Config.RateLimitRPS,
					Burst:      deps.Config.RateLimitBurst,
					TrustProxy: deps.Config.RateLimitTrustProxy,
				}))
			}
			ws.Handle("/ws", deps.WebSocket)
		})

		v1.Group(func(api chi.Router) {
			api.Use(chimw.Compress(compressionLevel))

			if deps.Config.RequestTimeout > 0 {
				api.Use(chimw.Timeout(deps.Config.RequestTimeout))
			}

			for _, limit := range rateLimiters(deps.Config) {
				api.Use(limit)
			}

			for _, module := range deps.Modules {
				module.RegisterRoutes(api)
			}
		})
	})

	return r
}

func metricsAuth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if token == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subtleCompare(r.Header.Get("Authorization"), "Bearer "+token) {
				next.ServeHTTP(w, r)

				return
			}

			http.NotFound(w, r)
		})
	}
}

func rateLimiters(cfg config.Config) []func(http.Handler) http.Handler {
	if !cfg.RateLimitEnabled {
		return nil
	}

	general := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:       cfg.RateLimitRPS,
		Burst:      cfg.RateLimitBurst,
		TrustProxy: cfg.RateLimitTrustProxy,
	})

	auth := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:       cfg.AuthRateLimitRPS,
		Burst:      cfg.AuthRateLimitBurst,
		TrustProxy: cfg.RateLimitTrustProxy,
	})

	return []func(http.Handler) http.Handler{general, onlyAuthRoutes(auth)}
}

func onlyAuthRoutes(limit func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		limited := limit(next)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, authPathPrefix) {
				limited.ServeHTTP(w, r)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func mountUploads(r chi.Router, publicURL, dir string, guard uploadGuard) {
	if publicURL == "" || dir == "" {
		return
	}

	prefix := "/" + strings.Trim(publicURL, "/")
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))

	r.Get(prefix+"/*", func(w http.ResponseWriter, req *http.Request) {
		displayID, name, ok := uploadPathParams(strings.TrimPrefix(req.URL.Path, prefix))
		if !ok {
			http.NotFound(w, req)

			return
		}

		itemID, visible, err := guard(req.Context(), displayID)
		if err != nil {
			writeStatus(w, http.StatusServiceUnavailable, `{"status":"unavailable"}`)

			return
		}

		if !visible {
			http.NotFound(w, req)

			return
		}

		w.Header().Set("Cache-Control", "private, max-age=86400, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		rewritten := req.Clone(req.Context())
		rewritten.URL.Path = prefix + "/" + itemID.String() + "/" + name
		fileServer.ServeHTTP(w, rewritten)
	})
}

func uploadPathParams(rest string) (string, string, bool) {
	parts := strings.Split(strings.Trim(path.Clean(rest), "/"), "/")
	if len(parts) != uploadPathParts {
		return "", "", false
	}

	if !domain.IsValidDisplayID(parts[0]) {
		return "", "", false
	}

	return parts[0], parts[1], true
}

func newUploadGuard(pool *pgxpool.Pool) uploadGuard {
	if pool == nil {
		return func(context.Context, string) (uuid.UUID, bool, error) { return uuid.Nil, false, nil }
	}

	return func(ctx context.Context, displayID string) (uuid.UUID, bool, error) {
		const query = `
			SELECT id FROM items
			WHERE display_id = $1 AND deleted_at IS NULL AND status IN ('published', 'sold')`

		var itemID uuid.UUID
		err := pool.QueryRow(ctx, query, displayID).Scan(&itemID)
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, false, nil
		}
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("check item visibility: %w", err)
		}

		return itemID, true, nil
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

func subtleCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func writeStatus(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(body)); err != nil {
		slog.Error("write response", slog.String("error", err.Error()))
	}
}
