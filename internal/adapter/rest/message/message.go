package message

import (
	"github.com/go-chi/chi/v5"
	"github.com/ziliscite/messaging-app/internal/core/service/message"
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
		r.Get("/history", h.GetAll)
	})
}
