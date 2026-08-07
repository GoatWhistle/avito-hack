package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ScanFunc[T any] func(pgx.Row) (T, error)

func QueryAll[T any](
	ctx context.Context,
	q Querier,
	capacity int,
	scan ScanFunc[T],
	query string,
	args ...any,
) ([]T, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}
	defer rows.Close()

	if capacity < 0 {
		capacity = 0
	}

	out := make([]T, 0, capacity)

	for rows.Next() {
		value, scanErr := scan(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		out = append(out, value)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return out, nil
}
