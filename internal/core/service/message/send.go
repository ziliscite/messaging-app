package message

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
	"time"
)

type SendRequest struct {
	From    string
	Message string
}

func (s *Service) Send(ctx context.Context, request *SendRequest) (*domain.Message, error) {
	message := &domain.Message{
		From:    request.From,
		Message: request.Message,
		Date:    time.Now(),
	}

	err := s.messageRepo.Insert(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
