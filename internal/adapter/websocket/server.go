package websocket

import (
	"context"
	"github.com/ziliscite/messaging-app/internal/core/service/message"
	"go.elastic.co/apm"
	"net/http"
)

func (s *Socket) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.up.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.logger.Printf("New connection from %s", r.RemoteAddr)

	s.add(conn)
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
		s.remove(conn)
	}()

	for {
		var msg message.SendRequest
		if err = conn.ReadJSON(&msg); err != nil {
			break
		}

		tx := apm.DefaultTracer.StartTransaction("send message", "websocket")
		ctx := apm.ContextWithTransaction(context.Background(), tx)

		response, err := s.service.Send(ctx, &msg)
		if err != nil {
			s.logger.Printf("Error sending message: %v", err)
		}

		s.broadcast <- *response
	}
}
