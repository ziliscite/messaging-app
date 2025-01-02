package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ziliscite/messaging-app/internal/config"
	"github.com/ziliscite/messaging-app/pkg"
	"net/http"
)

func main() {
	configs := config.New()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger, middleware.Recoverer, middleware.URLFormat, middleware.CleanPath)

	pkg.MustServe(run(configs, mux))
}

func run(configs *config.Config, mux *chi.Mux) error {
	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		//_ = json.NewEncoder(w).Encode("OK")
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "OK"})
	})

	fmt.Println("Environment: ", configs.Environment)
	fmt.Printf("Server is starting on %s\n", configs.Address())
	fmt.Println("Routes:")
	_ = chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	return http.ListenAndServe(configs.Address(), mux)
}
