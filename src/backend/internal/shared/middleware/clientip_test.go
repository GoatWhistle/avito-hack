package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/middleware"
)

func TestClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		trustProxy bool
		want       string
	}{
		{
			name:       "remote addr without proxy",
			remoteAddr: "203.0.113.5:4321",
			want:       "203.0.113.5",
		},
		{
			name:       "forwarded header ignored when proxy untrusted",
			remoteAddr: "203.0.113.5:4321",
			headers:    map[string]string{"X-Forwarded-For": "198.51.100.7"},
			want:       "203.0.113.5",
		},
		{
			name:       "forwarded for wins when trusted",
			remoteAddr: "10.0.0.1:1234",
			headers:    map[string]string{"X-Forwarded-For": "198.51.100.7, 10.0.0.1"},
			trustProxy: true,
			want:       "198.51.100.7",
		},
		{
			name:       "real ip fallback when trusted",
			remoteAddr: "10.0.0.1:1234",
			headers:    map[string]string{"X-Real-IP": "198.51.100.9"},
			trustProxy: true,
			want:       "198.51.100.9",
		},
		{
			name:       "garbage forwarded falls back to remote addr",
			remoteAddr: "203.0.113.5:4321",
			headers:    map[string]string{"X-Forwarded-For": "not-an-ip"},
			trustProxy: true,
			want:       "203.0.113.5",
		},
		{
			name:       "ipv6 remote addr",
			remoteAddr: "[2001:db8::1]:4321",
			want:       "2001:db8::1",
		},
		{
			name:       "remote addr without port",
			remoteAddr: "203.0.113.5",
			want:       "203.0.113.5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			req.RemoteAddr = tc.remoteAddr

			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			require.Equal(t, tc.want, middleware.ClientIP(req, tc.trustProxy))
		})
	}
}
