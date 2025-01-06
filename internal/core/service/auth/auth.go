package auth

import (
	"context"
	"github.com/ziliscite/messaging-app/config"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/session"
)

type Repository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
}

type API interface {
	CreateSession(ctx context.Context, request *SessionRequest) (*SessionResponse, error)
}

type Service struct {
	sessionRepo Repository
	tc          *config.TokenConfig
}

func New(repo Repository, secret *config.TokenConfig) *Service {
	return &Service{
		sessionRepo: repo,
		tc:          secret,
	}
}
