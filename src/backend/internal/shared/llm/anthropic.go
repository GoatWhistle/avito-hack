package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultModel     = "claude-haiku-4-5"
	DefaultEndpoint  = "https://api.anthropic.com/v1/messages"
	DefaultTimeout   = 5 * time.Second
	DefaultMaxTokens = 300
	apiVersion       = "2023-06-01"
	maxResponseBytes = 1 << 20
)

var (
	ErrNoAPIKey       = errors.New("anthropic api key is not configured")
	ErrEmptyResponse  = errors.New("anthropic returned an empty completion")
	ErrUnexpectedCode = errors.New("anthropic returned an unexpected status code")
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	APIKey    string
	Model     string
	Endpoint  string
	Timeout   time.Duration
	MaxTokens int
	HTTP      Doer
}

type AnthropicClient struct {
	apiKey    string
	model     string
	endpoint  string
	timeout   time.Duration
	maxTokens int
	http      Doer
}

func NewAnthropicClient(cfg Config) (*AnthropicClient, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, ErrNoAPIKey
	}

	client := &AnthropicClient{
		apiKey:    cfg.APIKey,
		model:     orDefault(cfg.Model, DefaultModel),
		endpoint:  orDefault(cfg.Endpoint, DefaultEndpoint),
		timeout:   cfg.Timeout,
		maxTokens: cfg.MaxTokens,
		http:      cfg.HTTP,
	}
	if client.timeout <= 0 {
		client.timeout = DefaultTimeout
	}
	if client.maxTokens <= 0 {
		client.maxTokens = DefaultMaxTokens
	}
	if client.http == nil {
		client.http = &http.Client{Timeout: client.timeout}
	}

	return client, nil
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []message `json:"messages"`
}

type responseBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type response struct {
	Content []responseBlock `json:"content"`
}

func (c *AnthropicClient) Complete(ctx context.Context, factsJSON string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	body, err := json.Marshal(request{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    SystemPrompt,
		Messages:  []message{{Role: "user", Content: fmt.Sprintf(UserPromptTemplate, factsJSON)}},
	})
	if err != nil {
		return "", fmt.Errorf("encode anthropic request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build anthropic request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	return c.send(req)
}

func (c *AnthropicClient) send(req *http.Request) (text string, err error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("call anthropic: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close anthropic response: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d", ErrUnexpectedCode, resp.StatusCode)
	}

	payload, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return "", fmt.Errorf("read anthropic response: %w", err)
	}

	var decoded response
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return "", fmt.Errorf("decode anthropic response: %w", err)
	}

	return textOf(decoded)
}

func textOf(decoded response) (string, error) {
	parts := make([]string, 0, len(decoded.Content))
	for _, block := range decoded.Content {
		if block.Type == "text" && strings.TrimSpace(block.Text) != "" {
			parts = append(parts, strings.TrimSpace(block.Text))
		}
	}

	if len(parts) == 0 {
		return "", ErrEmptyResponse
	}

	return strings.Join(parts, " "), nil
}

func orDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
