package auth

import (
	"context"
	"github.com/ziliscite/messaging-app/pkg/token"
	"go.elastic.co/apm"
	"time"
)

type RefreshRequest struct {
	UserID       uint   `json:"user_id,omitempty"`
	Email        string `json:"email,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshResponse struct {
	AccessToken          string    `json:"access_token,omitempty"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
}

func (s *Service) Refresh(ctx context.Context, refreshRequest *RefreshRequest) (*RefreshResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "refresh session", "service")
	defer span.End()

	session, err := s.sessionRepo.GetSession(spanCtx, refreshRequest.UserID, refreshRequest.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := token.Create(spanCtx, session.UserID, s.tc.AccessTokenExpirationTime, refreshRequest.Email, s.tc.Secret)
	if err != nil {
		return nil, err
	}

	err = s.sessionRepo.Update(spanCtx, accessToken, accessTokenExpiresAt, session.RefreshToken, session.UserID)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}, nil
}
