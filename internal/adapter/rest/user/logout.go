package user

import (
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/domain"
	"github.com/ziliscite/messaging-app/pkg/middleware"
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, fmt.Sprintf("Internal server error: %s", domain.ErrFailedParsingValue), http.StatusInternalServerError)
		return
	}

	err := h.authService.Revoke(r.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
