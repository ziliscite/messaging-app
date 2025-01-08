package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/service/auth"
	"github.com/ziliscite/messaging-app/internal/core/service/user"
	"github.com/ziliscite/messaging-app/pkg/res"
	"github.com/ziliscite/messaging-app/pkg/token"
	"net/http"
	"time"
)

type LoginResponse struct {
	ID                    uint      `json:"id,omitempty"`
	Username              string    `json:"username,omitempty"`
	Email                 string    `json:"email,omitempty"`
	AccessToken           string    `json:"access_token,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at,omitempty"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return session tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} res.BadRequestError "Bad Request - Invalid input data"
// @Failure 401 {object} res.UnauthorizedError "Unauthorized - Invalid credentials"
// @Failure 404 {object} res.NotFoundError "Not Found - User not found"
// @Failure 500 {object} res.InternalServerError "Internal Server Error"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var request user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		res.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userResponse, err := h.userService.Login(r.Context(), &request)
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrDatabase):
			res.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		case errors.Is(err, posgres.ErrNotFound):
			res.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, user.ErrInvalidCredentials):
			res.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			res.Error(w, err.Error(), http.StatusBadRequest)
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
			res.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		default:
			res.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	res.Success(w, LoginResponse{
		ID:                    userResponse.ID,
		Username:              userResponse.Username,
		Email:                 userResponse.Email,
		AccessToken:           sessionResponse.AccessToken,
		AccessTokenExpiresAt:  sessionResponse.AccessTokenExpiresAt,
		RefreshToken:          sessionResponse.RefreshToken,
		RefreshTokenExpiresAt: sessionResponse.RefreshTokenExpiresAt,
	}, http.StatusOK)
}
