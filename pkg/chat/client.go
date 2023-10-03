package chat

import (
	"bytes"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
	ID       uuid.UUID
	Username string

	connections map[uuid.UUID]*connection
}

type connection struct {
	room *room
	conn *websocket.Conn
	send chan []byte
}

func (c *client) readPump(roomId uuid.UUID) {
	connection := c.connections[roomId]

	defer func() {
		connection.room.unregister <- c
		connection.conn.Close()
		delete(c.connections, roomId)
	}()
	connection.conn.SetReadLimit(maxMessageSize)
	connection.conn.SetReadDeadline(time.Now().Add(pongWait))
	connection.conn.SetPongHandler(func(string) error { connection.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := connection.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		connection.room.broadcast <- message
	}
}

func (c *client) writePump(roomId uuid.UUID) {
	connection := c.connections[roomId]
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		connection.conn.Close()
		delete(c.connections, roomId)
	}()
	for {
		select {
		case message, ok := <-connection.send:
			connection.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room closed the channel.
				connection.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := connection.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(connection.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-connection.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			connection.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := connection.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
