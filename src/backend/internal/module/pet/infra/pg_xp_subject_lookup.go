package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

func (r *PgXPEventRepository) ExistsBySubject(
	ctx context.Context,
	userID uuid.UUID,
	action domain.Action,
	subjectID uuid.UUID,
) (bool, error) {
	const query = `SELECT EXISTS (
		SELECT 1 FROM xp_events WHERE user_id = $1 AND action = $2 AND subject_id = $3
	)`

	var exists bool
	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, string(action), subjectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check xp event subject: %w", err)
	}

	return exists, nil
}
