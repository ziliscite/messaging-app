package user

import (
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	err := h.authService.Revoke(r.Context())
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
