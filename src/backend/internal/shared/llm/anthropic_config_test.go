package llm_test

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/llm"
)

func TestNewAnthropicClientRejectsBlankKey(t *testing.T) {
	t.Parallel()

	cases := []string{"", "   ", "\t\n"}

	for _, key := range cases {
		t.Run("key="+strings.TrimSpace(key)+"|", func(t *testing.T) {
			t.Parallel()

			client, err := llm.NewAnthropicClient(llm.Config{APIKey: key})

			require.ErrorIs(t, err, llm.ErrNoAPIKey)
			assert.Nil(t, client)
		})
	}
}

func TestCompleteAppliesDefaultModelAndMaxTokens(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"ok"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{"total_xp":1}`)
	require.NoError(t, err)

	body := decodeRequest(t, doer)
	assert.Equal(t, llm.DefaultModel, body["model"])
	assert.InEpsilon(t, float64(llm.DefaultMaxTokens), body["max_tokens"], 1e-9)
	assert.Equal(t, llm.SystemPrompt, body["system"])
}

func TestCompleteHonoursExplicitConfig(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"ok"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{
		APIKey:    "key",
		Model:     "claude-custom",
		Endpoint:  "https://proxy.internal/v1/messages",
		MaxTokens: 42,
		HTTP:      doer,
	})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{}`)
	require.NoError(t, err)

	body := decodeRequest(t, doer)
	assert.Equal(t, "claude-custom", body["model"])
	assert.InEpsilon(t, 42.0, body["max_tokens"], 1e-9)
	assert.Equal(t, "https://proxy.internal/v1/messages", doer.captured.URL.String())
}

func TestCompleteFallsBackOnNonPositiveConfig(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"ok"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{
		APIKey: "key", HTTP: doer, Timeout: -time.Second, MaxTokens: -5, Model: "  ", Endpoint: " ",
	})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{}`)
	require.NoError(t, err)

	body := decodeRequest(t, doer)
	assert.Equal(t, llm.DefaultModel, body["model"])
	assert.InEpsilon(t, float64(llm.DefaultMaxTokens), body["max_tokens"], 1e-9)
	assert.Equal(t, llm.DefaultEndpoint, doer.captured.URL.String())
}

func TestCompleteRejectsUnbuildableEndpoint(t *testing.T) {
	t.Parallel()

	client, err := llm.NewAnthropicClient(llm.Config{
		APIKey: "key", Endpoint: "://not a url", HTTP: &stubDoer{},
	})
	require.NoError(t, err)

	text, err := client.Complete(t.Context(), `{}`)

	require.Error(t, err)
	assert.Empty(t, text)
}

func TestCompleteJoinsMultipleTextBlocks(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK,
		body: `{"content":[{"type":"text","text":" Первое "},{"type":"thinking","text":"skip"},` +
			`{"type":"text","text":"  "},{"type":"text","text":"Второе "}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	text, err := client.Complete(t.Context(), `{}`)

	require.NoError(t, err)
	assert.Equal(t, "Первое Второе", text)
}

func TestCompleteReportsUnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{status: http.StatusInternalServerError, body: `{}`}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{}`)

	require.ErrorIs(t, err, llm.ErrUnexpectedCode)
	assert.Contains(t, err.Error(), "500")
}

func TestCompleteReportsEmptyResponse(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{status: http.StatusOK, body: `{"content":[{"type":"text","text":"   "}]}`}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{}`)

	require.ErrorIs(t, err, llm.ErrEmptyResponse)
}

func TestCompletePropagatesBodyReadFailure(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{status: http.StatusOK, reader: failingReader{}}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{}`)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "read anthropic response")
}

func TestCompleteEmbedsFactsInUserMessage(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"ok"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	_, err = client.Complete(t.Context(), `{"total_xp":123}`)
	require.NoError(t, err)

	body := decodeRequest(t, doer)
	messages, ok := body["messages"].([]any)
	require.True(t, ok)
	require.Len(t, messages, 1)

	first, ok := messages[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "user", first["role"])
	assert.Contains(t, first["content"], `{"total_xp":123}`)
}

func TestGeneratorRejectsUnencodableFacts(t *testing.T) {
	t.Parallel()

	doer := &stubDoer{
		status: http.StatusOK, body: `{"content":[{"type":"text","text":"ok"}]}`,
	}
	client, err := llm.NewAnthropicClient(llm.Config{APIKey: "key", HTTP: doer})
	require.NoError(t, err)

	text, err := llm.NewGenerator(client, nil).
		GenerateSummary(t.Context(), map[string]float64{"nan": math.NaN()})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "encode summary facts")
	assert.Empty(t, text)
}

func TestGeneratorNilReceiverIsSafe(t *testing.T) {
	t.Parallel()

	var generator *llm.Generator

	text, err := generator.GenerateSummary(context.Background(), nil)

	require.ErrorIs(t, err, llm.ErrNoAPIKey)
	assert.Empty(t, text)
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, assert.AnError }

func decodeRequest(t *testing.T, doer *stubDoer) map[string]any {
	t.Helper()

	require.NotNil(t, doer.captured)
	raw, err := io.ReadAll(doer.captured.Body)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	return decoded
}
