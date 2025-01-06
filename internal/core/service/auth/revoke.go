package auth

import (
	"context"
	"errors"
	"github.com/ziliscite/messaging-app/pkg/middleware"
)

var ErrFailedParsingValue = errors.New("failed parsing value")

func (s *Service) Revoke(ctx context.Context) error {
	accessToken, ok := ctx.Value(middleware.AccessToken).(string)
	if !ok {
		return ErrFailedParsingValue
	}

	userId, ok := ctx.Value(middleware.UserIDKey).(uint)
	if !ok {
		return ErrFailedParsingValue
	}

	session, err := s.sessionRepo.GetSession(ctx, accessToken, userId)
	if err != nil {
		return err
	}

	return s.sessionRepo.Revoke(ctx, session.AccessToken, session.UserID)
}
