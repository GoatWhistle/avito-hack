package infra_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestEventDeduplicatorClaim(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tx      *pgtest.Tx
		want    bool
		wantErr error
	}{
		{name: "claimed", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}, want: true},
		{name: "already claimed", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}, want: false},
		{name: "error", tx: &pgtest.Tx{ExecErrs: []error{errDB}}, wantErr: errDB},
	}

	eventID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	event := events.Event{
		Type:       events.Type("item.published"),
		UserID:     userA,
		SubjectID:  itemA,
		OccurredAt: fixedTime,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := infra.NewPgEventDeduplicator(nil).Claim(ctxWith(tc.tx), eventID, event)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.False(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tc.tx.ExecCalls[0], "INSERT INTO pet_processed_events", "ON CONFLICT DO NOTHING")
			assert.Equal(t,
				[]any{eventID, "item.published", userA, itemA, fixedTime},
				tc.tx.ExecCalls[0].Args)
		})
	}
}

func TestEventDeduplicatorClaimTreatsMultiRowAsNotClaimed(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(2)}}

	got, err := infra.NewPgEventDeduplicator(nil).
		Claim(ctxWith(tx), uuid.New(), events.Event{UserID: userA, OccurredAt: fixedTime})

	require.NoError(t, err)
	assert.False(t, got)
}
