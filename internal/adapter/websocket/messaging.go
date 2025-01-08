package websocket

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/ziliscite/messaging-app/internal/core/domain/message"
	"github.com/ziliscite/messaging-app/pkg/must"
	"net/http"
	"sync"
	"time"
)

type Socket struct {
	mux       *chi.Mux
	up        *websocket.Upgrader
	clients   map[*websocket.Conn]bool
	m         sync.Mutex
	broadcast chan message.Message
	// service
}

func NewSocket(mux *chi.Mux) *Socket {
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
		broadcast: make(chan message.Message),
		m:         sync.Mutex{},
	}
}

func (s *Socket) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.up.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.add(conn)
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
		s.remove(conn)
	}()

	for {
		var msg message.Message
		if err = conn.ReadJSON(&msg); err != nil {
			break
		}

		msg.Date = time.Now()

		// insert

		s.broadcast <- msg
	}
}

func (s *Socket) HandleMessages() {
	for msg := range s.broadcast {
		clientsCopy := s.get()
		for client := range clientsCopy {
			if err := client.WriteJSON(msg); err != nil {
				err = client.Close()
				if err != nil {
					return
				}
				s.remove(client)
			}
		}
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
	s.mux.Get("/message/v1/send", s.HandleConnections)
	go s.HandleMessages()

	fmt.Printf("Starting WebSocket server on %s\n", address)
	must.MustServe(http.ListenAndServe(address, s.mux))
}
