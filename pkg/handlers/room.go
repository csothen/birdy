package handlers

import (
	"strconv"

	birdy "github.com/csothen/birdy/pkg/core"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  birdy.ChatService
	upgrader websocket.Upgrader
}

func NewHandler() *Handler {
	return &Handler{
		service:  birdy.NewService(),
		upgrader: websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096},
	}
}

func (h *Handler) JoinRoom(c echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	roomParam := c.Param("room")
	roomId, err := strconv.Atoi(roomParam)
	if err != nil {
		return err
	}

	room, err := h.service.JoinRoom(conn, roomId)
	if err != nil {
		return err
	}

	return c.Render(200, "room", room)
}
