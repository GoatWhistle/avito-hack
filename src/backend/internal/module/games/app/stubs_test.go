package app_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type stubClock struct {
	now time.Time
}

func (c stubClock) Now() time.Time { return c.now }

type stubTx struct {
	calls int
}

func (t *stubTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	t.calls++

	return fn(ctx)
}

type stubRounds struct {
	byDisplayID  map[string]*domain.Round
	active       *domain.Round
	saved        []*domain.Round
	saveErr      error
	lockedRounds []string
	lockedGames  []string
}

func (r *stubRounds) LockUserGame(_ context.Context, userID uuid.UUID, gameSlug string) error {
	r.lockedGames = append(r.lockedGames, userID.String()+"/"+gameSlug)

	return nil
}

func newStubRounds() *stubRounds {
	return &stubRounds{byDisplayID: map[string]*domain.Round{}}
}

func (r *stubRounds) Save(_ context.Context, round *domain.Round) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.byDisplayID[round.DisplayID()] = round
	r.saved = append(r.saved, round)

	return nil
}

func (r *stubRounds) ByDisplayID(_ context.Context, displayID string) (*domain.Round, error) {
	round, ok := r.byDisplayID[displayID]
	if !ok {
		return nil, domain.ErrRoundNotFound
	}

	return round, nil
}

func (r *stubRounds) ByDisplayIDForUpdate(ctx context.Context, displayID string) (*domain.Round, error) {
	r.lockedRounds = append(r.lockedRounds, displayID)

	return r.ByDisplayID(ctx, displayID)
}

func (r *stubRounds) ActiveByUser(_ context.Context, _ uuid.UUID, _ string) (*domain.Round, error) {
	if r.active == nil {
		return nil, domain.ErrRoundNotFound
	}

	return r.active, nil
}

type scoreKey struct {
	userID uuid.UUID
	slug   string
}

type stubScores struct {
	best  map[scoreKey]int
	saves []int
}

func newStubScores() *stubScores {
	return &stubScores{best: map[scoreKey]int{}}
}

func (s *stubScores) BestScore(_ context.Context, userID uuid.UUID, slug string) (int, error) {
	return s.best[scoreKey{userID: userID, slug: slug}], nil
}

func (s *stubScores) SaveBestScore(_ context.Context, userID uuid.UUID, slug string, score int) (int, error) {
	key := scoreKey{userID: userID, slug: slug}

	s.saves = append(s.saves, score)

	if score > s.best[key] {
		s.best[key] = score
	}

	return s.best[key], nil
}

type progressKey struct {
	userID uuid.UUID
	slug   string
	day    domain.Day
}

type dayKey struct {
	userID uuid.UUID
	day    domain.Day
}

type stubProgress struct {
	daily         map[progressKey]domain.DailyProgress
	dailyAny      map[dayKey]domain.DailyProgress
	streaks       map[uuid.UUID]domain.Streak
	lockedStreaks []string
}

func newStubProgress() *stubProgress {
	return &stubProgress{
		daily:    map[progressKey]domain.DailyProgress{},
		dailyAny: map[dayKey]domain.DailyProgress{},
		streaks:  map[uuid.UUID]domain.Streak{},
	}
}

func (p *stubProgress) Daily(
	_ context.Context,
	userID uuid.UUID,
	slug string,
	day domain.Day,
) (domain.DailyProgress, error) {
	return p.daily[progressKey{userID: userID, slug: slug, day: day}], nil
}

func (p *stubProgress) SaveDaily(
	_ context.Context,
	userID uuid.UUID,
	slug string,
	day domain.Day,
	progress domain.DailyProgress,
) error {
	p.daily[progressKey{userID: userID, slug: slug, day: day}] = progress

	return nil
}

func (p *stubProgress) DailyAny(
	_ context.Context,
	userID uuid.UUID,
	day domain.Day,
) (domain.DailyProgress, error) {
	return p.dailyAny[dayKey{userID: userID, day: day}], nil
}

func (p *stubProgress) SaveDailyAny(
	_ context.Context,
	userID uuid.UUID,
	day domain.Day,
	progress domain.DailyProgress,
) error {
	p.dailyAny[dayKey{userID: userID, day: day}] = progress

	return nil
}

func (p *stubProgress) Streak(_ context.Context, userID uuid.UUID) (domain.Streak, error) {
	return p.streaks[userID], nil
}

func (p *stubProgress) StreakForUpdate(ctx context.Context, userID uuid.UUID) (domain.Streak, error) {
	p.lockedStreaks = append(p.lockedStreaks, userID.String())

	return p.Streak(ctx, userID)
}

func (p *stubProgress) SaveStreak(_ context.Context, userID uuid.UUID, streak domain.Streak) error {
	p.streaks[userID] = streak

	return nil
}

type scriptedGame struct {
	slug     string
	target   int
	correct  bool
	progress domain.Progress
	next     bool
	err      error
}

func (g *scriptedGame) Slug() string      { return g.slug }
func (g *scriptedGame) TargetStreak() int { return g.target }

func (g *scriptedGame) Start(_ context.Context, r *domain.Round) (domain.View, error) {
	if g.err != nil {
		return domain.View{}, g.err
	}

	r.SetPayload(json.RawMessage(`{"step":1}`))

	return domain.View{Prompt: json.RawMessage(`{"question":"first"}`)}, nil
}

func (g *scriptedGame) Resume(_ context.Context, _ *domain.Round) (domain.View, error) {
	return domain.View{Prompt: json.RawMessage(`{"question":"resumed"}`)}, nil
}

func (g *scriptedGame) Guess(
	_ context.Context,
	_ *domain.Round,
	_ json.RawMessage,
) (domain.GuessOutcome, error) {
	if g.err != nil {
		return domain.GuessOutcome{}, g.err
	}

	outcome := domain.GuessOutcome{
		Correct:  g.correct,
		Progress: g.progress,
		Reveal:   json.RawMessage(`{"truth":42}`),
	}

	if g.correct || g.next {
		outcome.Next = &domain.View{Prompt: json.RawMessage(`{"question":"next"}`)}
	}

	return outcome, nil
}

type attemptGame struct {
	scriptedGame

	maxAttempts  int
	attemptsUsed int
}

func (g *attemptGame) MaxAttempts() int                 { return g.maxAttempts }
func (g *attemptGame) AttemptsUsed(_ *domain.Round) int { return g.attemptsUsed }
