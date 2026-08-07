package llm_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/llm"
)

type stubDoer struct {
	status   int
	body     string
	reader   io.Reader
	err      error
	captured *http.Request
}

func (s *stubDoer) Do(req *http.Request) (*http.Response, error) {
	s.captured = req
	if s.err != nil {
		return nil, s.err
	}

	body := s.reader
	if body == nil {
		body = strings.NewReader(s.body)
	}

	return &http.Response{
		StatusCode: s.status,
		Body:       io.NopCloser(body),
		Header:     make(http.Header),
	}, nil
}

func TestNewAnthropicClientRequiresKey(t *testing.T) {
	t.Parallel()

	client, err := llm.NewAnthropicClient(llm.Config{})

	require.ErrorIs(t, err, llm.ErrNoAPIKey)
	assert.Nil(t, client)
}

func TestCompleteReturnsText(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":" Сегодня был отличный день! "}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	text, err := client.Complete(context.Background(), `{"total_xp":7}`)

	require.NoError(t, err)
	assert.Equal(t, "Сегодня был отличный день!", text)
	require.NotNil(t, doer.captured)
	assert.Equal(t, "key", doer.captured.Header.Get("x-api-key"))
	assert.Equal(t, "2023-06-01", doer.captured.Header.Get("anthropic-version"))
}

func TestCompleteFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		doer *stubDoer
	}{
		{name: "transport error", doer: &stubDoer{err: errors.New("no network")}},
		{name: "non 200", doer: &stubDoer{status: http.StatusTooManyRequests, body: `{}`}},
		{name: "broken json", doer: &stubDoer{status: http.StatusOK, body: `not json`}},
		{name: "empty content", doer: &stubDoer{status: http.StatusOK, body: `{"content":[]}`}},
		{
			name: "only non text blocks",
			doer: &stubDoer{status: http.StatusOK, body: `{"content":[{"type":"thinking"}]}`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: tc.doer})
			require.NoError(t, err)

			text, err := client.Complete(context.Background(), `{}`)

			require.Error(t, err)
			assert.Empty(t, text)
		})
	}
}

func TestCompleteRespectsTimeout(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{err: context.DeadlineExceeded}
	client, err := llm.NewAnthropicClient(llm.Config{
		APIKey: "key", HTTP: doer, Timeout: 10 * time.Millisecond,
	})
	require.NoError(t, err)

	_, err = client.Complete(context.Background(), `{}`)

	require.Error(t, err)
}

func TestGeneratorFallsBackOnError(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{err: errors.New("dns failure")}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	text, err := llm.NewGenerator(client, nil).
		GenerateSummary(context.Background(), map[string]int{"total_xp": 3})

	require.Error(t, err)
	assert.Empty(t, text)
}

func TestGeneratorReturnsCompletion(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"Ура!"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	text, err := llm.NewGenerator(client, nil).
		GenerateSummary(context.Background(), map[string]int{"total_xp": 3})

	require.NoError(t, err)
	assert.Equal(t, "Ура!", text)
}

func TestGeneratorWithoutCompleter(t *testing.T) {
	t.Parallel()

	text, err := llm.NewGenerator(nil, nil).GenerateSummary(context.Background(), nil)

	require.ErrorIs(t, err, llm.ErrNoAPIKey)
	assert.Empty(t, text)
}

func TestPromptMentionsRules(t *testing.T) {
	t.Parallel()

	assert.Contains(t, llm.SystemPrompt, "от первого лица")
	assert.Contains(t, llm.SystemPrompt, "2-4 предложения")
	assert.Contains(t, llm.UserPromptTemplate, "%s")
}
