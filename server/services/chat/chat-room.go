package chat

import (
	"github.com/csothen/birdy/services/auth"
	"github.com/google/uuid"
)

type ChatRoom struct {
	ID      string
	Name    string
	Owner   *auth.User
	Clients map[*Client]struct{}

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewChatRoom(name string, owner *auth.User) *ChatRoom {
	return &ChatRoom{
		ID:      uuid.New().String(),
		Name:    name,
		Owner:   owner,
		Clients: make(map[*Client]struct{}),
	}
}

func (c *ChatRoom) Run() {
	select {
	case message := <-c.broadcast:
		c.sendMessage(message)
	case client := <-c.register:
		c.Clients[client] = struct{}{}
	case client := <-c.unregister:
		delete(c.Clients, client)
	}
}

func (c *ChatRoom) sendMessage(m []byte) {
	for client := range c.Clients {
		client.send <- m
	}
}
