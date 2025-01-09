package auth

import (
	"context"
	"go.elastic.co/apm"
)

func (s *Service) Revoke(ctx context.Context, userId uint) error {
	span, spanCtx := apm.StartSpan(ctx, "revoke session", "service")
	defer span.End()

	return s.sessionRepo.Revoke(spanCtx, userId)
}
