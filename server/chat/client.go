package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var newline = []byte{'\n'}

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type Client struct {
	Username string
	conn     *websocket.Conn
	send     chan []byte
	server   *Server
}

func NewClient(username string, conn *websocket.Conn, server *Server) *Client {
	return &Client{
		Username: username,
		conn:     conn,
		send:     make(chan []byte),
		server:   server,
	}
}

// readPump reads messages that the clients is trying to send
func (c *Client) readPump() {
	defer func() {
		c.disconnect()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.handleMessage(payload)
	}
}

// writePump writes the messages to the client's connection
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(payload []byte) {
	var m Message
	if err := json.Unmarshal(payload, &m); err != nil {
		log.Printf("Error on unmarshal message payload %s", err)
		return
	}

	m.Sender = c.Username
	c.server.handler <- m
}

func (c *Client) disconnect() {
	c.server.unregister <- c
}
