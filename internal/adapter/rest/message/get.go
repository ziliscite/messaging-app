package message

import (
	"github.com/ziliscite/messaging-app/pkg/res"
	"go.elastic.co/apm"
	"net/http"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	span, ctx := apm.StartSpan(r.Context(), "get history", "controller")
	defer span.End()

	histories, err := h.service.GetAll(ctx)
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, histories, http.StatusOK)
}
