package handlers

import (
	"github.com/csothen/birdy/pkg/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRoom(c echo.Context) error {
	val := c.Get(userContextKey)
	if val == nil {
		c.Logger().Infof("attempt to join room without token")
		return c.String(401, "unauthorized")
	}

	idParam := c.Param("id")
	roomId, err := uuid.ParseBytes([]byte(idParam))
	if err != nil {
		c.Logger().Errorf("error parsing ID param to UUID: %+v", err)
		return c.String(400, "room ID must be a valid UUID")
	}

	r, err := h.chatService.GetRoom(roomId)
	if err != nil {
		c.Logger().Errorf("error retrieving room with id '%s': %+v", roomId.String(), err)
		return c.String(500, "could not get room")
	}

	if r == nil {
		return c.String(404, "room not found")
	}

	page := room{
		ID:   r.ID.String(),
		Name: r.Name,
	}

	return c.Render(200, "room", page)
}

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
	roomId, err := uuid.ParseBytes([]byte(idParam))
	if err != nil {
		c.Logger().Errorf("error converting ID param to UUID: %+v", err)
		return c.String(400, "room ID must be a valid UUID")
	}

	if err := h.chatService.JoinRoom(conn, user.ID, roomId); err != nil {
		c.Logger().Errorf("error joining room: %+v", err)
		return c.String(500, "error joining room")
	}

	return nil
}
