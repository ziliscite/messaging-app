package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrDatabase       = errors.New("database error")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{
		db: conn,
	}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, username, email, created_at, updated_at
    `

	var createdUser domain.User
	if err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateEmail
		}

		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}

	return &createdUser, nil
}
