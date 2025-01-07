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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
