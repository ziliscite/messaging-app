package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziliscite/messaging-app/internal/adapter/posgres"
	"github.com/ziliscite/messaging-app/internal/core/service/user"
	"github.com/ziliscite/messaging-app/pkg/res"
	"go.elastic.co/apm"
	"net/http"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "User registration details"
// @Success 201 {object} user.RegisterResponse
// @Failure 400 {object} res.BadRequestError "Bad Request - Invalid input data"
// @Failure 409 {object} res.ConflictError "Conflict - User already exists"
// @Failure 500 {object} res.InternalServerError "Internal Server Error"
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	span, ctx := apm.StartSpan(r.Context(), "register", "controller")
	defer span.End()

	var request user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		res.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.userService.Register(ctx, &request)
	if err != nil {
		switch {
		case errors.Is(err, posgres.ErrDatabase) || errors.Is(err, user.ErrFailedHash):
			res.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		case errors.Is(err, posgres.ErrDuplicate):
			res.Error(w, err.Error(), http.StatusConflict)
		default:
			res.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	res.Success(w, response, http.StatusCreated)
}
