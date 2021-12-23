package chat

import (
	"github.com/google/uuid"
)

type Room struct {
	ID      string
	Name    string
	Owner   *Client
	Clients map[*Client]struct{}

	send       chan Message
	register   chan *Client
	unregister chan *Client
}

func NewRoom(name string, owner *Client) *Room {
	room := &Room{
		ID:      uuid.New().String(),
		Name:    name,
		Owner:   owner,
		Clients: make(map[*Client]struct{}),

		send:       make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	room.Clients[owner] = struct{}{}
	return room
}

func (r *Room) Run() {
	select {
	case message := <-r.send:
		r.sendMessage(message.encode())
	case client := <-r.register:
		r.handleJoin(client)
	case client := <-r.unregister:
		r.handleLeave(client)
	}
}

func (r *Room) sendMessage(m []byte) {
	for client := range r.Clients {
		client.send <- m
	}
}

func (r *Room) handleLeave(c *Client) {
	delete(r.Clients, c)
	if len(r.Clients) == 0 {
		c.server.handler <- Message{
			Sender:  r.Owner.Username,
			Type:    DeleteRoom,
			Content: r.ID,
		}
	}
}

func (r *Room) handleJoin(c *Client) {
	r.Clients[c] = struct{}{}
}
