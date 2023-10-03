package chat

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Service interface {
	ConnectClient(clientname string) error
	DisconnectClient(clientId uuid.UUID) error
	CreateRoom(clientId uuid.UUID, name string) (*room, error)
	DeleteRoom(clientId, roomId uuid.UUID) error
	JoinRoom(conn *websocket.Conn, clientId, roomId uuid.UUID) (*room, error)
	LeaveRoom(clientId, roomId uuid.UUID) error
}

type service struct {
	mutex   sync.Mutex
	clients map[uuid.UUID]*client
	rooms   map[uuid.UUID]*room
}

func NewService() Service {
	return &service{
		clients: make(map[uuid.UUID]*client),
		rooms:   make(map[uuid.UUID]*room),
	}
}

func (c *service) ConnectClient(clientname string) error {
	clientId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	client := &client{ID: clientId, Username: clientname, connections: make(map[uuid.UUID]*connection)}
	c.clients[clientId] = client

	return nil
}

func (c *service) DisconnectClient(clientId uuid.UUID) error {
	client, ok := c.clients[clientId]
	if !ok {
		return fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	for roomId := range client.connections {
		room := c.rooms[roomId]
		room.unregister <- client
	}
	delete(c.clients, clientId)

	return nil
}

func (c *service) CreateRoom(clientId uuid.UUID, name string) (*room, error) {
	owner, ok := c.clients[clientId]
	if !ok {
		return nil, fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	roomId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	room := &room{
		ID:         roomId,
		Name:       name,
		owner:      owner,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	c.rooms[roomId] = room
	go room.run()

	return room, nil
}

func (c *service) DeleteRoom(clientId, roomId uuid.UUID) error {
	room, ok := c.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%s' was not found", roomId.String())
	}

	if room.owner.ID != clientId {
		return fmt.Errorf("permission denied: not the owner")
	}

	for client := range room.clients {
		room.unregister <- client
	}

	delete(c.rooms, roomId)
	return nil
}

func (c *service) JoinRoom(conn *websocket.Conn, clientId, roomId uuid.UUID) (*room, error) {
	joiner, ok := c.clients[clientId]
	if !ok {
		return nil, fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	room, ok := c.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room with ID '%s' was not found", roomId.String())
	}

	roomConnection := &connection{room: room, conn: conn, send: make(chan []byte)}

	joiner.connections[roomId] = roomConnection

	roomConnection.room.register <- joiner
	go joiner.readPump(roomId)
	go joiner.writePump(roomId)

	return room, nil
}

func (c *service) LeaveRoom(clientId uuid.UUID, roomId uuid.UUID) error {
	client, ok := c.clients[clientId]
	if !ok {
		return fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	room, ok := c.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%s' does not exist", roomId.String())
	}

	room.unregister <- client
	delete(client.connections, roomId)
	return nil
}
