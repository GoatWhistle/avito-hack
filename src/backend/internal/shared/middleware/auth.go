package middleware

import (
	"net/http"
	"strings"

	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

const bearerPrefix = "Bearer "

type tokenParser interface {
	Parse(raw string) (auth.Actor, error)
}

func Authenticate(parser tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, err := extractBearer(r)
			if err != nil {
				apierr.Write(w, r, err)
				return
			}

			actor, err := parser.Parse(raw)
			if err != nil {
				apierr.Write(w, r, domainerr.ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), actor)))
		})
	}
}

func OptionalAuthenticate(parser tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, err := extractBearer(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			actor, err := parser.Parse(raw)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), actor)))
		})
	}
}

func RequireRole(roles ...auth.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor, err := auth.ActorFrom(r.Context())
			if err != nil {
				apierr.Write(w, r, err)
				return
			}

			if !actor.HasRole(roles...) {
				apierr.Write(w, r, domainerr.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractBearer(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", domainerr.ErrUnauthorized
	}

	if !strings.HasPrefix(header, bearerPrefix) {
		return "", domainerr.ErrUnauthorized
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if token == "" {
		return "", domainerr.ErrUnauthorized
	}

	return token, nil
}
