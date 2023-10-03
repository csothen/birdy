package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

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

	room, err := h.chatService.JoinRoom(conn, roomId)
	if err != nil {
		return err
	}

	return c.Render(200, "room", room)
}
