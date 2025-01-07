package user

import (
	"github.com/go-chi/chi/v5"
	"github.com/ziliscite/messaging-app/internal/core/service/auth"
	"github.com/ziliscite/messaging-app/internal/core/service/user"
	"github.com/ziliscite/messaging-app/pkg/middleware"
)

type Handler struct {
	userService user.API
	authService auth.API
}

func New(us user.API, as auth.API) *Handler {
	return &Handler{
		userService: us,
		authService: as,
	}
}

func (h *Handler) Routes(mux *chi.Mux) {
	mux.With().Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.With(middleware.AuthMiddleware).Delete("/logout", h.Logout)
		r.With(middleware.RefreshMiddleware).Put("/refresh", h.Refresh)
	})
}
