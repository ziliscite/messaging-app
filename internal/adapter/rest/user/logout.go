package user

import (
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/domain"
	"github.com/ziliscite/messaging-app/pkg/middleware"
	"github.com/ziliscite/messaging-app/pkg/res"
	"net/http"
)

// Logout godoc
// @Summary Logout user
// @Description Revoke the user's session
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer {token}"
// @Success 204 "No Content"
// @Failure 401 {object} res.UnauthorizedError "Unauthorized - Invalid or missing token"
// @Failure 404 {object} res.NotFoundError "Not Found - Session not found"
// @Failure 500 {object} res.InternalServerError "Internal Server Error"
// @Router /auth/logout [delete]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		res.Error(w, fmt.Sprintf("Internal server error: %s", domain.ErrFailedParsingValue), http.StatusInternalServerError)
		return
	}

	err := h.authService.Revoke(r.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrNotFound):
			res.Error(w, err.Error(), http.StatusNotFound)
		default:
			res.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	res.Success(w, nil, http.StatusNoContent)
}
