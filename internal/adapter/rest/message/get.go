package message

import (
	"github.com/ziliscite/messaging-app/pkg/res"
	"net/http"
)

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	histories, err := h.service.GetAll(r.Context())
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, histories, http.StatusOK)
}
