package session

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{
		db: conn,
	}
}
