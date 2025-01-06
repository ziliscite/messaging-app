package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/service/auth"
	"github.com/ziliscite/messaging-app/internal/core/service/user"
	"github.com/ziliscite/messaging-app/pkg/token"
	"net/http"
)

type LoginResponse struct {
	ID                    uint   `json:"id,omitempty"`
	Username              string `json:"username,omitempty"`
	Email                 string `json:"email,omitempty"`
	AccessToken           string `json:"access_token,omitempty"`
	AccessTokenExpiresAt  int64  `json:"access_token_expires_at,omitempty"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt int64  `json:"refresh_token_expires_at,omitempty"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userResponse, err := h.userService.Login(r.Context(), &request)
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrDatabase):
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		case errors.Is(err, posgres.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, user.ErrInvalidCredentials):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	sessionResponse, err := h.authService.CreateSession(r.Context(), &auth.SessionRequest{
		UserID: userResponse.ID,
		Email:  userResponse.Email,
	})
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrDatabase) || errors.Is(err, token.ErrFailedSigning):
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(LoginResponse{
		ID:                    userResponse.ID,
		Username:              userResponse.Username,
		Email:                 userResponse.Email,
		AccessToken:           sessionResponse.AccessToken,
		AccessTokenExpiresAt:  sessionResponse.AccessTokenExpiresAt.Unix(),
		RefreshToken:          sessionResponse.RefreshToken,
		RefreshTokenExpiresAt: sessionResponse.RefreshTokenExpiresAt.Unix(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
