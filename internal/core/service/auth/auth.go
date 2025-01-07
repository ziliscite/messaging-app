package auth

import (
	"context"
	"github.com/ziliscite/messaging-app/config"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/session"
	"time"
)

type Repository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
	GetSession(ctx context.Context, userId uint, refreshToken string) (*domain.Session, error)
	Revoke(ctx context.Context, userId uint) error
	Update(ctx context.Context, accessToken string, accessTokenExpiresAt time.Time, refreshToken string, userId uint) error
}

type API interface {
	CreateSession(ctx context.Context, request *SessionRequest) (*SessionResponse, error)
	Refresh(ctx context.Context, refreshRequest *RefreshRequest) (*RefreshResponse, error)
	Revoke(ctx context.Context, userId uint) error
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
