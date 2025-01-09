package websocket

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	domain "github.com/ziliscite/messaging-app/internal/core/domain/message"
	"github.com/ziliscite/messaging-app/internal/core/service/message"
	"github.com/ziliscite/messaging-app/pkg/must"
	"io"
	"log"
	"net/http"
	"sync"
)

type Socket struct {
	mux       *chi.Mux
	up        *websocket.Upgrader
	clients   map[*websocket.Conn]bool
	m         sync.Mutex
	broadcast chan domain.Message
	logger    *log.Logger
	service   message.WriteAPI
}

func NewSocket(mux *chi.Mux, wr io.Writer, service message.WriteAPI) *Socket {
	return &Socket{
		mux: mux,
		up: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:   make(map[*websocket.Conn]bool),
		m:         sync.Mutex{},
		broadcast: make(chan domain.Message),
		logger:    log.New(wr, "WebSocket: ", log.Ldate|log.Ltime),
		service:   service,
	}
}

func (s *Socket) add(conn *websocket.Conn) {
	s.m.Lock()
	defer s.m.Unlock()
	s.clients[conn] = true
}

func (s *Socket) remove(conn *websocket.Conn) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.clients, conn)
}

func (s *Socket) get() map[*websocket.Conn]bool {
	s.m.Lock()
	defer s.m.Unlock()
	clientsCopy := make(map[*websocket.Conn]bool, len(s.clients))
	for conn := range s.clients {
		clientsCopy[conn] = true
	}
	return clientsCopy
}

func (s *Socket) Start(address string) {
	s.mux.With().Route("/message", func(r chi.Router) {
		r.Get("/send", s.HandleConnections)
	})

	go s.HandleMessages()

	fmt.Printf("Starting WebSocket server on %s\n", address)
	must.MustServe(http.ListenAndServe(address, s.mux))
}
