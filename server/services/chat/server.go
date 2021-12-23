package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	authTimeout    = 1 * time.Minute
	authRetryLimit = 5
)

// Server represents the Chat Server
type Server struct {
	Clients map[string]*Client
	Rooms   map[string]*Room

	handler    chan Message
	register   chan *Client
	unregister chan *Client

	upgrader websocket.Upgrader
}

// NewServer creates a new Server
func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
		Rooms:   make(map[string]*Room),

		handler:    make(chan Message),
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
func (s *Server) Run() {
	select {
	case message := <-s.handler:
		s.handleMessage(message)
	case client := <-s.register:
		s.handleClientConnect(client)
	case client := <-s.unregister:
		s.handleClientDisconnect(client)
	}
}

// Start will take an address and start the Chat Server on that Address
func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		s.serve(rw, r)
	})

	log.Printf("server started on port %s\n", addr)

	return http.ListenAndServe(addr, nil)
}

// serve is responsible for taking the connection attempts
// and handle them by upgrading the request to a connection,
// authenticating the client and setting it up to read and write messages
func (s *Server) serve(rw http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), authTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- s.waitClientAuth(conn)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
}

func (s *Server) waitClientAuth(conn *websocket.Conn) error {
	retries := 0
	for {
		retries++
		if retries > authRetryLimit {
			return fmt.Errorf("authentication retry limit (%d) reached", authRetryLimit)
		}

		_, payload, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return fmt.Errorf("unexpected close error: %v", err)
			}
			return err
		}

		var message Message
		err = json.Unmarshal(payload, &message)
		if err != nil {
			return err
		}

		if message.Type != Authenticate {
			log.Printf("invalid message type when authenticating")
			continue
		}

		username := message.Content
		if _, ok := s.Clients[username]; ok {
			log.Printf("username %s already taken", username)
			continue
		}

		client := NewClient(username, conn, s)
		go client.writePump()
		go client.readPump()

		s.register <- client

		return nil
	}
}

func (s *Server) handleClientConnect(c *Client) {
	s.Clients[c.Username] = c
}

func (s *Server) handleClientDisconnect(c *Client) {
	delete(s.Clients, c.Username)
	for _, room := range s.Rooms {
		room.unregister <- c
	}
}

func (s *Server) handleMessage(m Message) {
	switch m.Type {
	case SendMessage:
		s.redirectMessage(m)
	case JoinRoom:
		s.addClientToRoom(m)
	case LeaveRoom:
		s.removeClientFromRoom(m)
	case CreateRoom:
		s.createRoom(m)
	case DeleteRoom:
		s.deleteRoom(m)
	default:
		log.Printf("server received invalid type of message")
	}
}

func (s *Server) redirectMessage(m Message) {
	if m.IsDM {
		client, ok := s.Clients[m.Target]
		if !ok {
			log.Printf("the target client (%s) does not exist", m.Target)
			return
		}
		client.send <- m.encode()
	} else {
		room, ok := s.Rooms[m.Target]
		if !ok {
			log.Printf("the target room (%s) does not exist", m.Target)
			return
		}
		room.send <- m
	}
}

func (s *Server) addClientToRoom(m Message) {
	sender, ok := s.Clients[m.Sender]
	if !ok {
		log.Printf("can't join room, sender (%s) does not exist", m.Sender)
		return
	}

	room, ok := s.Rooms[m.Target]
	if !ok {
		log.Printf("can't join room, room (%s) does not exist", m.Target)
		return
	}

	room.register <- sender
}

func (s *Server) removeClientFromRoom(m Message) {
	sender, ok := s.Clients[m.Sender]
	if !ok {
		log.Printf("can't leave room, sender (%s) does not exist", m.Sender)
		return
	}

	room, ok := s.Rooms[m.Target]
	if !ok {
		log.Printf("can't leave room, room (%s) does not exist", m.Target)
		return
	}

	room.unregister <- sender
}

func (s *Server) createRoom(m Message) {
	name := m.Target
	room := NewRoom(name, s.Clients[m.Sender])

	s.Rooms[room.ID] = room
}

func (s *Server) deleteRoom(m Message) {
	name := m.Target
	room, ok := s.Rooms[name]
	if !ok {
		log.Printf("can't delete room, room (%s) does not exist", name)
		return
	}

	if room.Owner.Username != m.Sender {
		log.Printf("can't delete room, user does not own the room")
		return
	}

	delete(s.Rooms, name)
}
