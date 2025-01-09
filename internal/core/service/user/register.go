package user

import (
	"context"
	"errors"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterResponse struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

var ErrFailedHash = errors.New("failed to hash password")

func (s *Service) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	span, spanCtx := apm.StartSpan(ctx, "register", "service")
	defer span.End()

	username, err := domain.ValidateUsername(request.Username)
	if err != nil {
		return nil, err
	}

	email, err := domain.ValidateEmail(request.Email)
	if err != nil {
		return nil, err
	}

	password, err := domain.ValidatePassword(request.Password)
	if err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedHash
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}

	user, err = s.userRepo.Create(spanCtx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
