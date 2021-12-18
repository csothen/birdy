package main

import (
	"log"

	"github.com/csothen/terminal-client/client"
	"github.com/csothen/terminal-client/ui"
	"github.com/gorilla/websocket"
)

func main() {
	establishConnection()
	ui.Render()
}

func establishConnection() {
	socketUrl := "ws://localhost:8080" + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	client := client.New(conn)
	client.Start()
}
