package message

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
)

type Repository interface {
	GetAll(ctx context.Context) (*[]domain.Message, error)
	Insert(ctx context.Context, message *domain.Message) error
}

type WriteAPI interface {
	Send(ctx context.Context, request *SendRequest) (*domain.Message, error)
}

type ReadAPI interface {
	GetAll(ctx context.Context) (*[]domain.Message, error)
}

type Service struct {
	messageRepo Repository
}

func New(repo Repository) *Service {
	return &Service{
		messageRepo: repo,
	}
}
