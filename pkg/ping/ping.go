package ping

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/ziliscite/messaging-app/pkg/res"
	"go.elastic.co/apm"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	// ctx dikirim ke turunan functionnya, macam service atau repository
	span, _ := apm.StartSpan(r.Context(), "ping", "controller")
	defer span.End()

	type PingRequest struct {
		Message string `json:"message"`
	}

	var request PingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		res.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Message) > 200 {
		res.Error(w, "message too long", http.StatusBadRequest)
		return
	}

	res.Success(w, map[string]any{"response": "Pong: " + request.Message}, http.StatusOK)
}

func health(w http.ResponseWriter, r *http.Request) {
	span, _ := apm.StartSpan(r.Context(), "health", "controller")
	defer span.End()
	res.Success(w, map[string]any{"message": "OK"}, http.StatusOK)
}

func Register(mux *chi.Mux) {
	mux.Group(func(r chi.Router) {
		mux.Get("/health", health)
		mux.Post("/ping", ping)
	})
}
