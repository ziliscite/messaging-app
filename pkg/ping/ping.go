package ping

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type PingRequest struct {
		Message string `json:"message"`
	}

	var request PingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Message) > 200 {
		http.Error(w, "message too long", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"response": "Pong: " + request.Message})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "OK"})
}

func Register(mux *chi.Mux) {
	mux.Group(func(r chi.Router) {
		mux.Get("/health", health)
		mux.Post("/ping", ping)
	})
}
