package birdy

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	id = 1
)

type ChatService interface {
	CreateRoom(name string) (*Room, error)
	DeleteRoom(roomId int) error
	JoinRoom(conn *websocket.Conn, roomId int) (*Room, error)
	LeaveRoom(roomId int) error
}

type chatService struct {
	mutex sync.Mutex
	rooms map[int]*Room
}

func NewService() ChatService {
	return &chatService{rooms: make(map[int]*Room)}
}

func (c *chatService) CreateRoom(name string) (*Room, error) {
	c.mutex.Lock()
	id++
	roomId := id
	c.mutex.Unlock()

	room := &Room{
		ID:         roomId,
		Name:       name,
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	c.rooms[roomId] = room
	go room.run()

	return room, nil
}

func (c *chatService) DeleteRoom(roomId int) error {
	delete(c.rooms, roomId)
	return nil
}

func (c *chatService) JoinRoom(conn *websocket.Conn, roomId int) (*Room, error) {
	room, ok := c.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room with ID '%d' does not exist", roomId)
	}

	joiner := &client{room: room, conn: conn, send: make(chan []byte)}

	joiner.room.register <- joiner
	go joiner.readPump()
	go joiner.writePump()

	return room, nil
}

func (c *chatService) LeaveRoom(roomId int) error {
	_, ok := c.rooms[roomId]
	if !ok {
		return fmt.Errorf("room with ID '%d' does not exist", roomId)
	}

	// we need to unregister the client, for that we need access to a client identifier
	return fmt.Errorf("not implemented")
}
