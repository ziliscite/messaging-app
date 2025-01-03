package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ziliscite/messaging-app/config"
	userRepository "github.com/ziliscite/messaging-app/internal/adapter/posgres/user"
	userHandler "github.com/ziliscite/messaging-app/internal/adapter/rest/user"
	userService "github.com/ziliscite/messaging-app/internal/core/service/user"
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

	UserMux(mux, conn)

	ping.Register(mux)

	fmt.Printf("Running on %s\n", configs.Address())
	fmt.Printf("Routes:\n")
	must.MustServe(chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("  %-7s %s\n", method, route)
		return nil
	}))

	must.MustServe(http.ListenAndServe(configs.Address(), mux))
}

func UserMux(mux *chi.Mux, conn *pgxpool.Pool) {
	repository := userRepository.New(conn)
	service := userService.New(repository)
	handler := userHandler.New(service)

	handler.Routes(mux)
}
