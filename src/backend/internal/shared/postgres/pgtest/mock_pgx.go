package pgtest

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNoStub = errors.New("pgtest: no stub configured")

type Call struct {
	SQL  string
	Args []any
}

type Row struct {
	Values []any
	Err    error
}

func (r Row) Scan(dest ...any) error {
	if r.Err != nil {
		return r.Err
	}

	return assign(dest, r.Values)
}

type Rows struct {
	Records [][]any
	IterErr error
	ScanErr error
	index   int
	Closed  bool
}

func (r *Rows) Close()                                       { r.Closed = true }
func (r *Rows) Err() error                                   { return r.IterErr }
func (r *Rows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *Rows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *Rows) Values() ([]any, error)                       { return nil, nil }
func (r *Rows) RawValues() [][]byte                          { return nil }
func (r *Rows) Conn() *pgx.Conn                              { return nil }

func (r *Rows) Next() bool {
	if r.index >= len(r.Records) {
		return false
	}
	r.index++

	return true
}

func (r *Rows) Scan(dest ...any) error {
	if r.ScanErr != nil {
		return r.ScanErr
	}

	return assign(dest, r.Records[r.index-1])
}

type Tx struct {
	QueryRows  []pgx.Rows
	QueryErrs  []error
	RowResults []Row
	ExecTags   []pgconn.CommandTag
	ExecErrs   []error

	QueryCalls    []Call
	QueryRowCalls []Call
	ExecCalls     []Call
	BatchCalls    []*pgx.Batch

	CommitErr   error
	RollbackErr error
	Committed   bool
	RolledBack  bool

	queryN    int
	queryRowN int
	execN     int
}

func (t *Tx) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	t.QueryCalls = append(t.QueryCalls, Call{SQL: sql, Args: args})

	if err := at(t.QueryErrs, t.queryN); err != nil {
		t.queryN++

		return nil, err
	}

	idx := t.queryN
	t.queryN++

	if idx < len(t.QueryRows) {
		return t.QueryRows[idx], nil
	}

	return &Rows{}, nil
}

func (t *Tx) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	t.QueryRowCalls = append(t.QueryRowCalls, Call{SQL: sql, Args: args})

	idx := t.queryRowN
	t.queryRowN++

	if idx < len(t.RowResults) {
		return t.RowResults[idx]
	}

	return Row{Err: ErrNoStub}
}

func (t *Tx) Exec(_ context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	t.ExecCalls = append(t.ExecCalls, Call{SQL: sql, Args: args})

	idx := t.execN
	t.execN++

	if err := at(t.ExecErrs, idx); err != nil {
		return pgconn.CommandTag{}, err
	}

	if idx < len(t.ExecTags) {
		return t.ExecTags[idx], nil
	}

	return pgconn.CommandTag{}, nil
}

func (t *Tx) SendBatch(_ context.Context, batch *pgx.Batch) pgx.BatchResults {
	t.BatchCalls = append(t.BatchCalls, batch)

	return nil
}

func (t *Tx) Begin(context.Context) (pgx.Tx, error) { return t, nil }

func (t *Tx) Commit(context.Context) error {
	t.Committed = true

	return t.CommitErr
}

func (t *Tx) Rollback(context.Context) error {
	t.RolledBack = true

	return t.RollbackErr
}

func (t *Tx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (t *Tx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, nil //nolint:nilnil // unused by the code under test
}

func (t *Tx) LargeObjects() pgx.LargeObjects { return pgx.LargeObjects{} }

func (t *Tx) Conn() *pgx.Conn { return nil }

func at(errs []error, i int) error {
	if i < len(errs) {
		return errs[i]
	}

	return nil
}

func Affected(n int64) pgconn.CommandTag {
	return pgconn.NewCommandTag("UPDATE " + itoa(n))
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}

	var buf [20]byte

	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}

	return string(buf[pos:])
}
