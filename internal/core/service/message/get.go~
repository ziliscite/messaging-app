package message

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
)

func (s *Service) GetAll(ctx context.Context) (*[]domain.Message, error) {
	return s.messageRepo.GetAll(ctx)
}
