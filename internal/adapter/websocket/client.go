package websocket

func (s *Socket) HandleMessages() {
	for msg := range s.broadcast {
		s.logger.Printf("Broadcasting message: %v", msg)
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
