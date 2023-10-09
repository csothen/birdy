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
	ListRooms() []*room
	GetRoom(id uuid.UUID) (*room, error)
	CreateRoom(clientId uuid.UUID, name string) (*room, error)
	DeleteRoom(clientId, roomId uuid.UUID) error
	JoinRoom(conn *websocket.Conn, clientId, roomId uuid.UUID) error
	LeaveRoom(clientId, roomId uuid.UUID) error
}

type service struct {
	mutex   sync.Mutex
	clients map[uuid.UUID]*client
	rooms   map[uuid.UUID]*room
}

func NewService() Service {
	r1Id, _ := uuid.NewUUID()
	r2Id, _ := uuid.NewUUID()
	r3Id, _ := uuid.NewUUID()
	r4Id, _ := uuid.NewUUID()

	r1 := &room{
		ID:         r1Id,
		Name:       "Room 1",
		owner:      nil,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	r2 := &room{
		ID:         r2Id,
		Name:       "Room 2",
		owner:      nil,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	r3 := &room{
		ID:         r3Id,
		Name:       "Room 3",
		owner:      nil,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	r4 := &room{
		ID:         r4Id,
		Name:       "Room 4",
		owner:      nil,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	rooms := map[uuid.UUID]*room{
		r1.ID: r1,
		r2.ID: r2,
		r3.ID: r3,
		r4.ID: r4,
	}

	return &service{
		clients: make(map[uuid.UUID]*client),
		rooms:   rooms,
	}
}

func (s *service) ConnectClient(clientname string) error {
	clientId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	client := &client{ID: clientId, Username: clientname, connections: make(map[uuid.UUID]*connection)}
	s.clients[clientId] = client

	return nil
}

func (s *service) DisconnectClient(clientId uuid.UUID) error {
	client, ok := s.clients[clientId]
	if !ok {
		return fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	for roomId := range client.connections {
		room := s.rooms[roomId]
		room.unregister <- client
	}
	delete(s.clients, clientId)

	return nil
}

func (s *service) ListRooms() []*room {
	rooms := []*room{}
	for _, r := range s.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

func (s *service) GetRoom(id uuid.UUID) (*room, error) {
	room, ok := s.rooms[id]
	if !ok {
		return nil, nil
	}
	return room, nil
}

func (s *service) CreateRoom(clientId uuid.UUID, name string) (*room, error) {
	owner, ok := s.clients[clientId]
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

	s.rooms[roomId] = room
	go room.run()

	return room, nil
}

func (s *service) DeleteRoom(clientId, roomId uuid.UUID) error {
	room, ok := s.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%s' was not found", roomId.String())
	}

	if room.owner.ID != clientId {
		return fmt.Errorf("permission denied: not the owner")
	}

	for client := range room.clients {
		room.unregister <- client
	}

	delete(s.rooms, roomId)
	return nil
}

func (s *service) JoinRoom(conn *websocket.Conn, clientId, roomId uuid.UUID) error {
	joiner, ok := s.clients[clientId]
	if !ok {
		return fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	room, ok := s.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%s' was not found", roomId.String())
	}

	roomConnection := &connection{room: room, conn: conn, send: make(chan []byte)}

	joiner.connections[roomId] = roomConnection

	roomConnection.room.register <- joiner
	go joiner.readPump(roomId)
	go joiner.writePump(roomId)

	return nil
}

func (s *service) LeaveRoom(clientId uuid.UUID, roomId uuid.UUID) error {
	client, ok := s.clients[clientId]
	if !ok {
		return fmt.Errorf("client with ID '%s' was not found", clientId.String())
	}

	room, ok := s.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%s' does not exist", roomId.String())
	}

	room.unregister <- client
	delete(client.connections, roomId)
	return nil
}
