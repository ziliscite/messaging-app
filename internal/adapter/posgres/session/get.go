package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/session"
)

// GetSession returns a refresh token by access token and user id
func (r *Repository) GetSession(ctx context.Context, accessToken string, userId uint) (*domain.Session, error) {
	query := `
		SELECT id, user_id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, created_at, updated_at
		FROM sessions
		WHERE refresh_token = $1 AND user_id = $2
	`

	var sess domain.Session
	if err := r.db.QueryRow(ctx, query, accessToken, userId).Scan(
		&sess.ID,
		&sess.UserID,
		&sess.AccessToken,
		&sess.AccessTokenExpiresAt,
		&sess.RefreshToken,
		&sess.RefreshTokenExpiresAt,
		&sess.CreatedAt,
		&sess.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session %w for user %d", posgres.ErrNotFound, userId)
		}

		return nil, fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	return &sess, nil
}
