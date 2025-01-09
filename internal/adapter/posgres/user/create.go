package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
	"go.elastic.co/apm"
)

func (r *Repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	span, spanCtx := apm.StartSpan(ctx, "create", "repository")
	defer span.End()

	query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at, updated_at
    `

	var createdUser domain.User
	if err := r.db.QueryRow(spanCtx, query, user.Username, user.Email, user.Password).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, posgres.ErrDuplicate
		}

		return nil, fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	return &createdUser, nil
}
