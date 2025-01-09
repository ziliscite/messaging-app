package message

import (
	"github.com/go-chi/chi"
	"github.com/ziliscite/messaging-app/internal/core/service/message"
	"github.com/ziliscite/messaging-app/pkg/middleware"
)

type Handler struct {
	service message.ReadAPI
}

func New(service message.ReadAPI) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes(mux *chi.Mux) {
	mux.With().Route("/message", func(r chi.Router) {
		r.With(middleware.AuthMiddleware).Get("/history", h.GetAll)
	})
}
