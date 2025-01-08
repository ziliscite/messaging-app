package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaggo/http-swagger/v2"
	"github.com/ziliscite/messaging-app/config"
	_ "github.com/ziliscite/messaging-app/docs"
	sessionRepository "github.com/ziliscite/messaging-app/internal/adapter/posgres/session"
	userRepository "github.com/ziliscite/messaging-app/internal/adapter/posgres/user"
	userHandler "github.com/ziliscite/messaging-app/internal/adapter/rest/user"
	authService "github.com/ziliscite/messaging-app/internal/core/service/auth"
	userService "github.com/ziliscite/messaging-app/internal/core/service/user"
	"github.com/ziliscite/messaging-app/pkg/db/posgres"
	"github.com/ziliscite/messaging-app/pkg/must"
	"github.com/ziliscite/messaging-app/pkg/ping"
	"html/template"
	"net/http"
)

// @title Messaging App
// @version 1.0
// @description This is a simple messaging app API.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @host localhost:3000
// @BasePath /
func main() {
	configs := config.New()

	conn := posgres.New(configs.Database)
	defer conn.Close()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger, middleware.Recoverer, middleware.URLFormat, middleware.CleanPath)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	ping.Register(mux)
	UserMux(mux, configs.Token, conn)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
	})

	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	fmt.Printf("Running on %s\n", configs.Address())
	fmt.Printf("Routes:\n")
	must.MustServe(chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("  %-7s %s\n", method, route)
		return nil
	}))

	must.MustServe(http.ListenAndServe(configs.Address(), mux))
}

func UserMux(mux *chi.Mux, cfg *config.TokenConfig, conn *pgxpool.Pool) {
	userRepo := userRepository.New(conn)
	userSvc := userService.New(userRepo)

	sessionRepo := sessionRepository.New(conn)
	authSvc := authService.New(sessionRepo, cfg)

	handler := userHandler.New(userSvc, authSvc)
	handler.Routes(mux)
}

// swag init -g cmd/web/main.go
