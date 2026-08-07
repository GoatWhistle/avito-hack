package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const summaryColumns = `id, user_id, date, facts, message, advice, generated_by, created_at`

type PgSummaryRepository struct {
	pool *pgxpool.Pool
}

func NewPgSummaryRepository(pool *pgxpool.Pool) *PgSummaryRepository {
	return &PgSummaryRepository{pool: pool}
}

func (r *PgSummaryRepository) ByUserAndDate(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) (*domain.DailySummary, error) {
	const query = `SELECT ` + summaryColumns +
		` FROM daily_summaries WHERE user_id = $1 AND date = $2`

	row := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID, date)

	summary, err := scanSummary(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainerr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (r *PgSummaryRepository) Insert(ctx context.Context, summary *domain.DailySummary) (bool, error) {
	const query = `INSERT INTO daily_summaries (` + summaryColumns + `)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (user_id, date) DO NOTHING`

	facts, err := json.Marshal(summary.Facts())
	if err != nil {
		return false, fmt.Errorf("encode summary facts: %w", err)
	}

	advice, err := encodeAdvice(summary.Advice())
	if err != nil {
		return false, err
	}

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		summary.ID(), summary.UserID(), summary.Date(), facts, summary.Message(),
		advice, string(summary.GeneratedBy()), summary.CreatedAt(),
	)
	if err != nil {
		return false, fmt.Errorf("insert daily summary: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *PgSummaryRepository) ListByUser(
	ctx context.Context,
	f app.SummaryHistoryFilter,
) ([]*domain.DailySummary, error) {
	query := `SELECT ` + summaryColumns + ` FROM daily_summaries WHERE user_id = $1`
	args := []any{f.UserID}

	if !f.Before.IsZero() {
		query += ` AND date < $2`
		args = append(args, f.Before)
	}
	query += fmt.Sprintf(` ORDER BY date DESC, id DESC LIMIT $%d`, len(args)+1)
	args = append(args, f.Limit)

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query daily summaries: %w", err)
	}
	defer rows.Close()

	summaries := make([]*domain.DailySummary, 0, f.Limit)
	for rows.Next() {
		summary, scanErr := scanSummary(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		summaries = append(summaries, summary)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate daily summaries: %w", err)
	}

	return summaries, nil
}

type summaryScanner interface {
	Scan(dest ...any) error
}

func scanSummary(row summaryScanner) (*domain.DailySummary, error) {
	var (
		params      domain.RestoreSummaryParams
		rawFacts    []byte
		rawAdvice   []byte
		generatedBy string
	)

	err := row.Scan(&params.ID, &params.UserID, &params.Date, &rawFacts,
		&params.Message, &rawAdvice, &generatedBy, &params.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(rawFacts, &params.Facts); err != nil {
		return nil, fmt.Errorf("decode summary facts: %w", err)
	}
	if len(rawAdvice) > 0 {
		var advice domain.Advice
		if err := json.Unmarshal(rawAdvice, &advice); err != nil {
			return nil, fmt.Errorf("decode summary advice: %w", err)
		}
		params.Advice = &advice
	}

	params.GeneratedBy = domain.SummarySource(generatedBy)

	return domain.RestoreDailySummary(params), nil
}

func encodeAdvice(advice *domain.Advice) ([]byte, error) {
	if advice == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(advice)
	if err != nil {
		return nil, fmt.Errorf("encode summary advice: %w", err)
	}

	return encoded, nil
}
