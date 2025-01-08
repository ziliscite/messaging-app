package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/domain"
	"github.com/ziliscite/messaging-app/internal/core/service/auth"
	"github.com/ziliscite/messaging-app/pkg/middleware"
	"net/http"
)

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh the user's access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} auth.RefreshResponse
// @Failure 400 {object} res.BadRequestError "Bad Request - Invalid input data"
// @Failure 401 {object} res.UnauthorizedError "Unauthorized - Invalid or expired refresh token"
// @Failure 500 {object} res.InternalServerError "Internal Server Error"
// @Router /auth/refresh [put]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, fmt.Sprintf("Internal server error: %s", domain.ErrFailedParsingValue), http.StatusInternalServerError)
		return
	}

	email, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok {
		http.Error(w, fmt.Sprintf("Internal server error: %s", domain.ErrFailedParsingValue), http.StatusInternalServerError)
		return
	}

	refreshToken, ok := r.Context().Value(middleware.RefreshKey).(string)
	if !ok {
		http.Error(w, fmt.Sprintf("Internal server error: %s", domain.ErrFailedParsingValue), http.StatusInternalServerError)
		return
	}

	response, err := h.authService.Refresh(r.Context(), &auth.RefreshRequest{
		UserID:       userId,
		Email:        email,
		RefreshToken: refreshToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrNotFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
