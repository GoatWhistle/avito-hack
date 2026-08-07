package server_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/config"
	"github.com/avito-hack/backend/internal/server"
)

func freeAddr(t *testing.T) string {
	t.Helper()

	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := ln.Addr().String()
	require.NoError(t, ln.Close())

	return addr
}

func waitForServer(t *testing.T, addr string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	dialer := net.Dialer{Timeout: 100 * time.Millisecond}

	for time.Now().Before(deadline) {
		conn, err := dialer.DialContext(t.Context(), "tcp", addr)
		if err == nil {
			require.NoError(t, conn.Close())

			return
		}
	}

	t.Fatalf("server at %s did not become reachable", addr)
}

func testConfig(addr string) config.Config {
	return config.Config{
		HTTPAddr:          addr,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ShutdownTimeout:   5 * time.Second,
	}
}

func TestNewBuildsServer(t *testing.T) {
	t.Parallel()

	srv := server.New(testConfig(":8080"), http.NotFoundHandler(), discardLogger())

	assert.NotNil(t, srv)
}

func TestServerRunServesRequestsThenShutsDownOnContextCancel(t *testing.T) {
	t.Parallel()

	addr := freeAddr(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive"))
	})

	srv := server.New(testConfig(addr), handler, discardLogger())

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() { errCh <- srv.Run(ctx) }()

	waitForServer(t, addr)

	resp, err := http.Get("http://" + addr + "/ping") //nolint:noctx // short-lived probe against a local test server
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "alive", string(body))

	cancel()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestServerRunReturnsListenError(t *testing.T) {
	t.Parallel()

	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() { _ = ln.Close() })

	srv := server.New(testConfig(ln.Addr().String()), http.NotFoundHandler(), discardLogger())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Run(ctx) }()

	select {
	case runErr := <-errCh:
		require.Error(t, runErr)
		assert.Contains(t, runErr.Error(), "listen and serve")
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not report the listen failure")
	}
}

func TestServerRunShutsDownWithAlreadyCancelledContext(t *testing.T) {
	t.Parallel()

	srv := server.New(testConfig(freeAddr(t)), http.NotFoundHandler(), discardLogger())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Run(ctx) }()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return for a pre-cancelled context")
	}
}
