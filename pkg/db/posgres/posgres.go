package posgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ziliscite/messaging-app/config"
	"github.com/ziliscite/messaging-app/pkg/must"
)

func New(configs *config.Config) *pgxpool.Pool {
	return must.Must(pgxpool.New(context.Background(), configs.Database.ConnectionString()))
}
