package pgtest

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type BatchResults struct {
	Tags     []pgconn.CommandTag
	ExecErrs []error
	CloseErr error
	Closed   bool

	execN int
}

func (b *BatchResults) Exec() (pgconn.CommandTag, error) {
	idx := b.execN
	b.execN++

	if err := at(b.ExecErrs, idx); err != nil {
		return pgconn.CommandTag{}, err
	}

	if idx < len(b.Tags) {
		return b.Tags[idx], nil
	}

	return pgconn.CommandTag{}, nil
}

func (b *BatchResults) Query() (pgx.Rows, error) { return &Rows{}, nil }

func (b *BatchResults) QueryRow() pgx.Row { return Row{Err: ErrNoStub} }

func (b *BatchResults) Close() error {
	b.Closed = true

	return b.CloseErr
}

type BatchTx struct {
	Tx

	Results *BatchResults
}

func (t *BatchTx) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	t.Tx.SendBatch(ctx, batch)

	if t.Results == nil {
		t.Results = &BatchResults{}
	}

	return t.Results
}
