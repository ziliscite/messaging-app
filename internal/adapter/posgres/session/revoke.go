package session

import (
	"context"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
)

func (r *Repository) Revoke(ctx context.Context, userId uint) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userId)
	if err != nil {
		return fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session %w for user %d", posgres.ErrNotFound, userId)
	}

	return nil
}
