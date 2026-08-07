//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8080"
	dialTimeout    = 2 * time.Second
	requestTimeout = 15 * time.Second
)

type client struct {
	t       *testing.T
	baseURL string
	http    *http.Client
	token   string
}

type errorEnvelope struct {
	Error struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		Field     string `json:"field"`
		RequestID string `json:"request_id"`
	} `json:"error"`
}

func baseURL() string {
	if raw := strings.TrimSpace(os.Getenv("INTEGRATION_BASE_URL")); raw != "" {
		return strings.TrimRight(raw, "/")
	}

	return defaultBaseURL
}

func requireStack(t *testing.T) *client {
	t.Helper()

	if os.Getenv("INTEGRATION") == "" && os.Getenv("INTEGRATION_BASE_URL") == "" {
		t.Skip("set INTEGRATION=1 (or INTEGRATION_BASE_URL) to run integration tests against a running stack")
	}

	url := baseURL()
	if !reachable(url) {
		t.Skipf("stack is not reachable at %s, start it with: make up", url)
	}

	c := &client{
		t:       t,
		baseURL: url,
		http:    &http.Client{Timeout: requestTimeout},
	}
	c.waitReady()

	return c
}

func reachable(rawURL string) bool {
	host := strings.TrimPrefix(strings.TrimPrefix(rawURL, "http://"), "https://")
	host = strings.SplitN(host, "/", 2)[0]

	if !strings.Contains(host, ":") {
		host += ":80"
	}

	conn, err := net.DialTimeout("tcp", host, dialTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()

	return true
}

func (c *client) waitReady() {
	c.t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := c.http.Get(c.baseURL + "/readyz")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	c.t.Skipf("stack at %s did not become ready in time", c.baseURL)
}

func (c *client) do(method, path string, body any) (int, []byte) {
	c.t.Helper()

	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			c.t.Fatalf("marshal %s %s: %v", method, path, err)
		}
		reader = bytes.NewReader(encoded)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		c.t.Fatalf("build %s %s: %v", method, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		c.t.Fatalf("request %s %s: %v", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		c.t.Fatalf("read %s %s: %v", method, path, err)
	}

	return resp.StatusCode, payload
}

func (c *client) mustDo(method, path string, body any, want int, out any) {
	c.t.Helper()

	status, payload := c.do(method, path, body)
	if status != want {
		c.t.Fatalf("%s %s: got %d want %d, body: %s", method, path, status, want, truncate(payload))
	}
	if out == nil {
		return
	}
	if err := json.Unmarshal(payload, out); err != nil {
		c.t.Fatalf("decode %s %s: %v, body: %s", method, path, err, truncate(payload))
	}
}

func decodeError(t *testing.T, payload []byte) errorEnvelope {
	t.Helper()

	var envelope errorEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("decode error envelope: %v, body: %s", err, truncate(payload))
	}

	return envelope
}

func truncate(payload []byte) string {
	const limit = 512
	if len(payload) <= limit {
		return string(payload)
	}

	return string(payload[:limit]) + "..."
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@integration.test", prefix, time.Now().UnixNano())
}
