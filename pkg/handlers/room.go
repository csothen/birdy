package handlers

import (
	"github.com/csothen/birdy/pkg/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) JoinRoom(c echo.Context) error {
	val := c.Get(userContextKey)
	if val == nil {
		c.Logger().Infof("attempt to join room without token")
		return c.String(401, "unauthorized")
	}

	user, ok := val.(*auth.User)
	if !ok {
		c.Logger().Errorf("could not convert value in context to user")
		return c.String(500, "unknown error")
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Errorf("error upgrading to websocket: %+v", err)
		return c.String(500, "unknown error")
	}

	idParam := c.Param("id")
	roomId, err := uuid.FromBytes([]byte(idParam))
	if err != nil {
		c.Logger().Errorf("error converting ID param to int: %+v", err)
		return c.String(400, "room ID must be a string")
	}

	room, err := h.chatService.JoinRoom(conn, user.ID, roomId)
	if err != nil {
		c.Logger().Errorf("error joining room: %+v", err)
		return c.String(500, "error joining room")
	}

	return c.Render(200, "room", room)
}
