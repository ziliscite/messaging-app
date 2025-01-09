package session

import (
	"context"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"go.elastic.co/apm"
	"time"
)

func (r *Repository) Update(ctx context.Context, accessToken string, accessTokenExpiresAt time.Time, refreshToken string, userId uint) error {
	span, spanCtx := apm.StartSpan(ctx, "update session", "repository")
	defer span.End()

	query := `
		UPDATE sessions 
		SET access_token = $1, access_token_expires_at = $2 
		WHERE refresh_token = $3 AND user_id = $4
	`

	tag, err := r.db.Exec(spanCtx, query, accessToken, accessTokenExpiresAt, refreshToken, userId)
	if err != nil {
		return fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session %w for user %d", posgres.ErrNotFound, userId)
	}

	return nil
}
