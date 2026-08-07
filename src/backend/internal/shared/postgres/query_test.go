package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/postgres"
)

var errBoom = errors.New("boom")

type fakeRows struct {
	values   []int
	index    int
	iterErr  error
	closed   bool
	scanFail bool
}

func (r *fakeRows) Close()                                       { r.closed = true }
func (r *fakeRows) Err() error                                   { return r.iterErr }
func (r *fakeRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *fakeRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *fakeRows) Values() ([]any, error)                       { return nil, nil }
func (r *fakeRows) RawValues() [][]byte                          { return nil }
func (r *fakeRows) Conn() *pgx.Conn                              { return nil }

func (r *fakeRows) Next() bool {
	if r.index >= len(r.values) {
		return false
	}
	r.index++

	return true
}

func (r *fakeRows) Scan(dest ...any) error {
	if r.scanFail {
		return errBoom
	}

	target, ok := dest[0].(*int)
	if !ok {
		return errBoom
	}
	*target = r.values[r.index-1]

	return nil
}

type fakeQuerier struct {
	rows     *fakeRows
	queryErr error
	gotQuery string
	gotArgs  []any
}

func (q *fakeQuerier) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	q.gotQuery, q.gotArgs = sql, args

	if q.queryErr != nil {
		return nil, q.queryErr
	}

	return q.rows, nil
}

func (q *fakeQuerier) QueryRow(context.Context, string, ...any) pgx.Row { return nil }

func (q *fakeQuerier) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { return nil }

func (q *fakeQuerier) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func scanInt(row pgx.Row) (int, error) {
	var v int
	if err := row.Scan(&v); err != nil {
		return 0, err
	}

	return v, nil
}

func TestQueryAllCollectsRowsAndPassesArgs(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{rows: &fakeRows{values: []int{1, 2, 3}}}

	got, err := postgres.QueryAll(context.Background(), q, 3, scanInt, "SELECT n", 42)

	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, got)
	assert.Equal(t, "SELECT n", q.gotQuery)
	assert.Equal(t, []any{42}, q.gotArgs)
	assert.True(t, q.rows.closed)
}

func TestQueryAllReturnsEmptySliceNotNil(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{rows: &fakeRows{}}

	got, err := postgres.QueryAll(context.Background(), q, 0, scanInt, "SELECT n")

	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestQueryAllWrapsQueryError(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{queryErr: errBoom}

	_, err := postgres.QueryAll(context.Background(), q, 0, scanInt, "SELECT n")

	require.ErrorIs(t, err, errBoom)
}

func TestQueryAllPropagatesScanError(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{rows: &fakeRows{values: []int{1}, scanFail: true}}

	_, err := postgres.QueryAll(context.Background(), q, 1, scanInt, "SELECT n")

	require.ErrorIs(t, err, errBoom)
	assert.True(t, q.rows.closed)
}

func TestQueryAllWrapsIterationError(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{rows: &fakeRows{values: []int{1}, iterErr: errBoom}}

	_, err := postgres.QueryAll(context.Background(), q, 1, scanInt, "SELECT n")

	require.ErrorIs(t, err, errBoom)
}

func TestQueryAllHandlesNegativeCapacity(t *testing.T) {
	t.Parallel()

	q := &fakeQuerier{rows: &fakeRows{values: []int{7}}}

	got, err := postgres.QueryAll(context.Background(), q, -5, scanInt, "SELECT n")

	require.NoError(t, err)
	assert.Equal(t, []int{7}, got)
}
