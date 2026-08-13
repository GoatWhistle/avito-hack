package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type GuessCommand struct {
	UserID    uuid.UUID
	GameSlug  string
	RoundID   string
	Move      json.RawMessage
	ClientDay domain.Day
}

type GuessHandler struct {
	registry *domain.Registry
	rounds   RoundRepository
	progress ProgressRepository
	scores   ScoreRepository
	tx       TxManager
	clock    Clock
}

func NewGuessHandler(
	registry *domain.Registry,
	rounds RoundRepository,
	progress ProgressRepository,
	scores ScoreRepository,
	tx TxManager,
	clock Clock,
) *GuessHandler {
	return &GuessHandler{
		registry: registry,
		rounds:   rounds,
		progress: progress,
		scores:   scores,
		tx:       tx,
		clock:    clock,
	}
}

func (h *GuessHandler) Handle(ctx context.Context, cmd GuessCommand) (GuessResult, error) {
	game, err := h.registry.Get(cmd.GameSlug)
	if err != nil {
		return GuessResult{}, err
	}

	var result GuessResult

	err = h.tx.WithTx(ctx, func(ctx context.Context) error {
		round, loadErr := h.load(ctx, cmd, game.Slug())
		if loadErr != nil {
			return loadErr
		}

		outcome, guessErr := game.Guess(ctx, round, cmd.Move)
		if guessErr != nil {
			return guessErr
		}

		applyProgress(round, outcome, game.TargetStreak(), h.clock.Now())

		if saveErr := h.rounds.Save(ctx, round); saveErr != nil {
			return saveErr
		}

		result = GuessResult{
			Correct:      outcome.Correct,
			Progress:     progressOf(outcome),
			Reveal:       outcome.Reveal,
			Streak:       round.Streak(),
			AttemptsUsed: domain.AttemptsUsedOf(game, round),
			MaxAttempts:  domain.MaxAttemptsOf(game),
			State:        round.State(),
		}

		if outcome.Next != nil && round.IsActive() {
			result.Prompt = outcome.Next.Prompt
		}

		if round.State() != domain.StateWon {
			return nil
		}

		scored, isScored := domain.ScoredOf(game)

		if isScored {
			best, scoreErr := h.persistScore(ctx, cmd, scored.RoundScore(round))
			if scoreErr != nil {
				return scoreErr
			}

			result.BestScore = best
			result.Reveal = withPersistedBest(result.Reveal, scored.RoundScore(round), best)

			if !scored.CountsTowardStreak(scored.RoundScore(round)) {
				return nil
			}
		}

		result.AttemptCompleted = true

		streak, progressErr := h.completeAttempt(ctx, cmd, round)
		if progressErr != nil {
			return progressErr
		}

		view := toStreakView(streak)
		result.StreakAfter = &view

		return nil
	})
	if err != nil {
		return GuessResult{}, err
	}

	return result, nil
}

func applyProgress(round *domain.Round, outcome domain.GuessOutcome, targetStreak int, now time.Time) {
	switch progressOf(outcome) {
	case domain.ProgressContinue:
		round.Touch(now)
	case domain.ProgressAdvance:
		round.Advance(targetStreak, now)
	case domain.ProgressWin:
		round.Win(now)
	case domain.ProgressLose:
		round.Lose(now)
	}
}

func progressOf(outcome domain.GuessOutcome) domain.Progress {
	if outcome.Progress != "" {
		return outcome.Progress
	}

	if outcome.Correct {
		return domain.ProgressAdvance
	}

	return domain.ProgressLose
}

func (h *GuessHandler) load(ctx context.Context, cmd GuessCommand, slug string) (*domain.Round, error) {
	round, err := h.rounds.ByDisplayIDForUpdate(ctx, cmd.RoundID)
	if err != nil {
		return nil, err
	}

	if round.UserID() != cmd.UserID || round.GameSlug() != slug {
		return nil, domain.ErrRoundNotFound
	}

	if !round.IsActive() {
		return nil, domain.ErrRoundFinished(round.State())
	}

	return round, nil
}

func withPersistedBest(reveal json.RawMessage, score, best int) json.RawMessage {
	if len(reveal) == 0 {
		return reveal
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(reveal, &payload); err != nil {
		return reveal
	}

	if _, ok := payload["best_score"]; !ok {
		return reveal
	}

	bestJSON, err := json.Marshal(best)
	if err != nil {
		return reveal
	}

	newBestJSON, err := json.Marshal(score >= best && score > 0)
	if err != nil {
		return reveal
	}

	payload["best_score"] = bestJSON
	payload["new_best"] = newBestJSON

	patched, err := json.Marshal(payload)
	if err != nil {
		return reveal
	}

	return patched
}

func (h *GuessHandler) persistScore(ctx context.Context, cmd GuessCommand, score int) (int, error) {
	if h.scores == nil {
		return 0, nil
	}

	best, err := h.scores.SaveBestScore(ctx, cmd.UserID, cmd.GameSlug, score)
	if err != nil {
		return 0, err
	}

	return best, nil
}

func (h *GuessHandler) completeAttempt(
	ctx context.Context,
	cmd GuessCommand,
	round *domain.Round,
) (domain.Streak, error) {
	today := cmd.ClientDay

	daily, err := h.progress.Daily(ctx, cmd.UserID, cmd.GameSlug, today)
	if err != nil {
		return domain.Streak{}, err
	}

	if err := h.progress.SaveDaily(
		ctx, cmd.UserID, cmd.GameSlug, today, domain.AdvanceDaily(daily, round.BestStreak()),
	); err != nil {
		return domain.Streak{}, err
	}

	anyDaily, err := h.progress.DailyAny(ctx, cmd.UserID, today)
	if err != nil {
		return domain.Streak{}, err
	}

	firstOfDay := domain.IsFirstAttemptOfDay(anyDaily)

	if err := h.progress.SaveDailyAny(
		ctx, cmd.UserID, today, domain.AdvanceDaily(anyDaily, round.BestStreak()),
	); err != nil {
		return domain.Streak{}, err
	}

	streak, err := h.progress.StreakForUpdate(ctx, cmd.UserID)
	if err != nil {
		return domain.Streak{}, err
	}

	if !firstOfDay {
		return streak, nil
	}

	streak = domain.AdvanceStreak(streak, today)

	if err := h.progress.SaveStreak(ctx, cmd.UserID, streak); err != nil {
		return domain.Streak{}, err
	}

	return streak, nil
}
