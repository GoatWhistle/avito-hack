package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/shared/middleware"
)

const readyPingTimeout = 2 * time.Second

type ModuleRegistrar interface {
	RegisterRoutes(r chi.Router)
}

type RouterDeps struct {
	Config  config.Config
	Pool    *pgxpool.Pool
	Logger  *slog.Logger
	Modules []ModuleRegistrar
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(middleware.Recoverer(deps.Logger))
	r.Use(middleware.AccessLog(deps.Logger))
	r.Use(middleware.Metrics())
	r.Use(chimw.Timeout(deps.Config.RequestTimeout))
	r.Use(middleware.CORS(deps.Config.AllowedOrigins))

	r.Get("/healthz", healthHandler)
	r.Get("/readyz", readyHandler(deps.Pool))
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(v1 chi.Router) {
		for _, module := range deps.Modules {
			module.RegisterRoutes(v1)
		}
	})

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeStatus(w, http.StatusOK, `{"status":"ok"}`)
}

func readyHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
