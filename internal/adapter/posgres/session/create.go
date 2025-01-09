package session

import (
	"context"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/session"
	"go.elastic.co/apm"
)

func (r *Repository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	span, spanCtx := apm.StartSpan(ctx, "create session", "repository")
	defer span.End()

	query := `
		INSERT INTO sessions (user_id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, user_id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at, created_at, updated_at
	`

	var createdSession domain.Session
	if err := r.db.QueryRow(spanCtx, query, session.UserID, session.AccessToken, session.AccessTokenExpiresAt, session.RefreshToken, session.RefreshTokenExpiresAt).Scan(
		&createdSession.ID,
		&createdSession.UserID,
		&createdSession.AccessToken,
		&createdSession.AccessTokenExpiresAt,
		&createdSession.RefreshToken,
		&createdSession.RefreshTokenExpiresAt,
		&createdSession.CreatedAt,
		&createdSession.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	return &createdSession, nil
}
