package client

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	done      chan interface{}
	interrupt chan os.Signal
}

func New(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		done:      make(chan interface{}),
		interrupt: make(chan os.Signal),
	}
}

func (c *Client) handle() {
	defer close(c.done)
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		c.handlePayload(payload)
	}
}

func (c *Client) Start() {
	signal.Notify(c.interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	go c.handle()

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(`{ "action": 0, "content": "Hello from GolangDocs!"}`))
			if err != nil {
				log.Println("Error during writing to websocket:", err)
				return
			}

		case <-c.interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-c.done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}
