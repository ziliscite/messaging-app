package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
	"go.elastic.co/apm"
)

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	span, spanCtx := apm.StartSpan(ctx, "get by email", "repository")
	defer span.End()

	query := `
		SELECT id, username, password, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	if err := r.db.QueryRow(spanCtx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user", posgres.ErrNotFound)
		}

		return nil, fmt.Errorf("%w: %v", posgres.ErrDatabase, err)
	}

	return &user, nil
}
