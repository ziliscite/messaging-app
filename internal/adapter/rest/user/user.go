package user

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	repo "github.com/ziliscite/messaging-app/internal/adapter/posgres/user"
	"github.com/ziliscite/messaging-app/internal/core/service/user"
	"net/http"
)

type Handler struct {
	s user.API
}

func New(service user.API) *Handler {
	return &Handler{
		s: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var request user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.s.Register(r.Context(), &request)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDatabase) || errors.Is(err, user.ErrFailedHash):
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) Routes(mux *chi.Mux) {
	mux.With().Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		//r.Post("/login", h.Login)
	})
}
