package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaggo/http-swagger/v2"
	"github.com/ziliscite/messaging-app/config"
	"github.com/ziliscite/messaging-app/docs"
	messageRepository "github.com/ziliscite/messaging-app/internal/adapter/mongo/message"
	sessionRepository "github.com/ziliscite/messaging-app/internal/adapter/posgres/session"
	userRepository "github.com/ziliscite/messaging-app/internal/adapter/posgres/user"
	messageHandler "github.com/ziliscite/messaging-app/internal/adapter/rest/message"
	userHandler "github.com/ziliscite/messaging-app/internal/adapter/rest/user"
	"github.com/ziliscite/messaging-app/internal/adapter/websocket"
	authService "github.com/ziliscite/messaging-app/internal/core/service/auth"
	messageService "github.com/ziliscite/messaging-app/internal/core/service/message"
	userService "github.com/ziliscite/messaging-app/internal/core/service/user"
	logfile "github.com/ziliscite/messaging-app/logs"
	nsql "github.com/ziliscite/messaging-app/pkg/db/mongo"
	"github.com/ziliscite/messaging-app/pkg/db/posgres"
	cm "github.com/ziliscite/messaging-app/pkg/middleware"
	"github.com/ziliscite/messaging-app/pkg/must"
	"github.com/ziliscite/messaging-app/pkg/ping"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmchi"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"io"
	"log"
	"net/http"
)

// @title Messaging App
// @version 1.0
// @description This is a simple messaging app API.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	wr := logfile.Set()
	log.SetOutput(wr)

	configs := config.New()

	conn := posgres.New(configs.Database)
	defer conn.Close()

	mux := chi.NewRouter()
	mux.Use(cm.CustomLogger, middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	client := nsql.New(configs.Mongo)
	defer client.Disconnect(context.Background())

	// separate webserver connection for websockets
	socketMux := chi.NewRouter()
	socketMux.Use(cm.CustomLogger, middleware.Recoverer)

	// apm middleware
	tracer, err := apm.NewTracer("messaging-app", "1.0.0")
	if err != nil {
		log.Printf("Failed to initialize APM: %v", err)
	} else {
		mux.Use(apmchi.Middleware(apmchi.WithTracer(tracer)))
		socketMux.Use(apmchi.Middleware(apmchi.WithTracer(tracer)))
	}

	socket := MessageMux(mux, wr, socketMux, client)
	go socket.Start(configs.WebsocketAddress())

	ping.Register(mux)
	UserMux(mux, configs.Token, conn)
	Statics(mux, configs.Address())
	
	Serve(mux, configs.Address())
}

func MessageMux(mux *chi.Mux, wr io.Writer, socketMux *chi.Mux, client *mongo.Client) *websocket.Socket {
	messageRepo := messageRepository.New(client)
	messageSvc := messageService.New(messageRepo)

	handler := messageHandler.New(messageSvc)
	handler.Routes(mux)

	return websocket.NewSocket(socketMux, wr, messageSvc)
}

func UserMux(mux *chi.Mux, cfg *config.TokenConfig, conn *pgxpool.Pool) {
	userRepo := userRepository.New(conn)
	userSvc := userService.New(userRepo)

	sessionRepo := sessionRepository.New(conn)
	authSvc := authService.New(sessionRepo, cfg)

	handler := userHandler.New(userSvc, authSvc)
	handler.Routes(mux)
}

func Statics(mux *chi.Mux, address string) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, nil)
	})

	docs.SwaggerInfo.Host = address
	docs.SwaggerInfo.BasePath = "/"
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
}

func Serve(mux *chi.Mux, address string) {
	fmt.Printf("Webserver running on %s\n", address)
	fmt.Printf("Routes:\n")
	must.MustServe(chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("  %-7s %s\n", method, route)
		return nil
	}))

	must.MustServe(http.ListenAndServe(address, mux))
}

// swag init -g cmd/web/main.go
