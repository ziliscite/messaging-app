package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ziliscite/messaging-app/internal/config"
	"github.com/ziliscite/messaging-app/internal/util"
	"net/http"
)

func main() {
	configs := config.New()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger, middleware.Recoverer, middleware.URLFormat, middleware.CleanPath)

	util.MustServe(run(configs, mux))
}

func run(configs *config.Config, mux *chi.Mux) error {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		//_ = json.NewEncoder(w).Encode("OK")
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "OK"})
	})

	mux.Post("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type PingRequest struct {
			Message string `json:"message"`
		}

		var ping PingRequest
		if err := json.NewDecoder(r.Body).Decode(&ping); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if len(ping.Message) > 200 {
			http.Error(w, "message too long", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"response": "Pong: " + ping.Message})
	})

	fmt.Println("=== Server Initialization ===")
	fmt.Printf("Environment: %s\n", configs.Environment)
	fmt.Printf("Address: %s\n", configs.Address())

	fmt.Println("\n=== Registered Routes ===")
	err := chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("  %-7s %s\n", method, route)
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("\n=== Server Starting ===")

	return http.ListenAndServe(configs.Address(), mux)
}
