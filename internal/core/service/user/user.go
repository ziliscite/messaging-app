package user

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type API interface {
	Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
}

type Service struct {
	userRepo Repository
}

func New(userRepository Repository) *Service {
	return &Service{
		userRepo: userRepository,
	}
}
