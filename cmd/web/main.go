package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ziliscite/messaging-app/config"
	"github.com/ziliscite/messaging-app/pkg/db/posgres"
	"github.com/ziliscite/messaging-app/pkg/must"
	"github.com/ziliscite/messaging-app/pkg/ping"
	"net/http"
)

func main() {
	configs := config.New()

	conn := posgres.New(configs)
	defer conn.Close()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger, middleware.Recoverer, middleware.URLFormat, middleware.CleanPath)

	must.MustServe(run(configs, mux))
}

func run(configs *config.Config, mux *chi.Mux) error {
	mux.Get("/health", ping.Health)
	mux.Post("/ping", ping.Ping)

	if err := chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("  %-7s %s\n", method, route)
		return nil
	}); err != nil {
		return err
	}

	return http.ListenAndServe(configs.Address(), mux)
}
