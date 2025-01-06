package user

import (
	"context"
	"errors"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *Service) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	// These 2 should probably goes in the client
	_, err := domain.ValidateEmail(request.Email)
	if err != nil {
		return nil, err
	}

	_, err = domain.ValidatePassword(request.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
