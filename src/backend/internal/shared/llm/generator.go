package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Completer interface {
	Complete(ctx context.Context, factsJSON string) (string, error)
}

type SummaryGenerator interface {
	GenerateSummary(ctx context.Context, facts any) (string, error)
}

type Generator struct {
	completer Completer
	logger    *slog.Logger
}

func NewGenerator(completer Completer, logger *slog.Logger) *Generator {
	if logger == nil {
		logger = slog.Default()
	}

	return &Generator{completer: completer, logger: logger}
}

func (g *Generator) GenerateSummary(ctx context.Context, facts any) (string, error) {
	if g == nil || g.completer == nil {
		return "", ErrNoAPIKey
	}

	encoded, err := json.Marshal(facts)
	if err != nil {
		return "", fmt.Errorf("encode summary facts: %w", err)
	}

	text, err := g.completer.Complete(ctx, string(encoded))
	if err != nil {
		g.logger.WarnContext(ctx, "llm summary generation failed, falling back to template",
			slog.Any("error", err))

		return "", err
	}

	return text, nil
}
