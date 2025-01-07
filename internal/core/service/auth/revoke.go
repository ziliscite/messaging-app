package auth

import (
	"context"
)

func (s *Service) Revoke(ctx context.Context, userId uint) error {
	return s.sessionRepo.Revoke(ctx, userId)
}
