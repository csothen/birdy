package handlers

import (
	"github.com/csothen/birdy/pkg/auth"
	"github.com/csothen/birdy/pkg/chat"
	"github.com/gorilla/websocket"
)

type Handler struct {
	authService auth.Service
	chatService chat.Service
	upgrader    websocket.Upgrader
}

func NewHandler() *Handler {
	return &Handler{
		authService: auth.NewService(),
		chatService: chat.NewService(),
		upgrader:    websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096},
	}
}
