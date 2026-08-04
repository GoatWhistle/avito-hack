package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const healthcheckTimeout = 3 * time.Second

func isHealthcheckMode() bool {
	return len(os.Args) > 1 && os.Args[1] == "healthcheck"
}

func runHealthcheck(addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), healthcheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthcheckURL(addr), http.NoBody)
	if err != nil {
		return fmt.Errorf("healthcheck request: %w", err)
	}

	client := &http.Client{Timeout: healthcheckTimeout}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("healthcheck request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "healthcheck: close body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthcheck status: %d", resp.StatusCode)
	}

	return nil
}

func healthcheckURL(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "http://127.0.0.1:8080/healthz"
	}

	if host == "" || strings.EqualFold(host, "0.0.0.0") {
		host = "127.0.0.1"
	}

	return fmt.Sprintf("http://%s/healthz", net.JoinHostPort(host, port))
}
