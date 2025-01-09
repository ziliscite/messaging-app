package auth

import (
	"context"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/session"
	"github.com/ziliscite/messaging-app/pkg/token"
	"go.elastic.co/apm"
	"time"
)

type SessionRequest struct {
	UserID uint   `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

type SessionResponse struct {
	UserID                uint      `json:"user_id,omitempty"`
	AccessToken           string    `json:"access_token,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at,omitempty"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

func (s *Service) CreateSession(ctx context.Context, request *SessionRequest) (*SessionResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "create session", "service")
	defer span.End()

	refreshToken, refreshTokenExpiresAt, err := token.Create(spanCtx, request.UserID, s.tc.RefreshTokenExpirationTime, request.Email, s.tc.Secret)
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := token.Create(spanCtx, request.UserID, s.tc.AccessTokenExpirationTime, request.Email, s.tc.Secret)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionRepo.Create(spanCtx,
		&domain.Session{
			UserID:                request.UserID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessTokenExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshTokenExpiresAt,
		},
	)
	if err != nil {
		return nil, err
	}

	return &SessionResponse{
		UserID:                session.UserID,
		AccessToken:           session.AccessToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshToken:          session.RefreshToken,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
	}, nil
}
