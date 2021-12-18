package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server represents the Interface that a Chat Server should comply with
type Server interface {
	// Start will take an address and start the Chat Server on that Address
	Start(addr string) error

	// Run is responsible for listening for new messages that
	// are sent across the connections established
	Run()
}

// server is the implementation of the Server interface and represents
// the actual Chat Server
type server struct {
	clients map[*Client]struct{}
	rooms   map[string]*ChatRoom

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	upgrader websocket.Upgrader
}

// NewServer creates a new Server
func NewServer() Server {
	return &server{
		clients:    make(map[*Client]struct{}),
		rooms:      make(map[string]*ChatRoom),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
		},
	}
}

// Run is responsible for listening for new messages that
// are sent across the connections established
func (s *server) Run() {
	select {
	case content := <-s.broadcast:
		s.notify(content)

	case client := <-s.register:
		s.clients[client] = struct{}{}

	case client := <-s.unregister:
		delete(s.clients, client)
	}
}

func (s *server) notify(c []byte) {
	notification := &Message{
		Action:  Notify,
		Content: string(c),
	}
	for c := range s.clients {
		c.send <- notification.encode()
	}
}

// Start will take an address and start the Chat Server on that Address
func (s *server) Start(addr string) error {
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		s.serve(rw, r)
	})

	fmt.Printf("server started on port %s\n", addr)

	return http.ListenAndServe(addr, nil)
}

// serve is responsible for taking the connection attempts
// and handle them by upgrading the request to a connection
// and setting up the client to read and write messages
func (s *server) serve(rw http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(nil, conn, s)

	go client.writePump()
	go client.readPump()

	s.register <- client
}
