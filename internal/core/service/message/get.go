package message

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
	"go.elastic.co/apm"
)

func (s *Service) GetAll(ctx context.Context) (*[]domain.Message, error) {
	span, spanCtx := apm.StartSpan(ctx, "get history", "service")
	defer span.End()

	return s.messageRepo.GetAll(spanCtx)
}
