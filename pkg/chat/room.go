package chat

import "github.com/google/uuid"

type room struct {
	ID   uuid.UUID
	Name string

	owner      *client
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func (r *room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				close(client.connections[r.ID].send)
				delete(client.connections, r.ID)
				delete(r.clients, client)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.connections[r.ID].send <- message:
				default:
					close(client.connections[r.ID].send)
					delete(client.connections, r.ID)
					delete(r.clients, client)
				}
			}
		}
	}
}
