package handlers

import (
	"net/http"

	"github.com/csothen/birdy/pkg/pages"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (h *Handler) Login(c echo.Context) error {
	var lr LoginRequest

	if err := c.Bind(&lr); err != nil {
		c.Logger().Errorf("could not bind payload: %+v", err)
		return c.String(400, "invalid payload")
	}

	user, token, err := h.authService.Authenticate(lr.Username, lr.Password)
	if err != nil {
		c.Logger().Warnf("could not authenticate user: %+v", err)
		return c.String(401, "invalid username and password")
	}

	c.SetCookie(&http.Cookie{
		Name:    authCookieName,
		Value:   token.Value,
		Expires: token.Expiration,
	})

	rooms := h.chatService.ListRooms()
	roomsList := []pages.RoomData{}
	for _, r := range rooms {
		roomsList = append(roomsList, pages.RoomData{
			ID:   r.ID.String(),
			Name: r.Name,
		})
	}

	page := pages.IndexPage{
		Head: pages.HeadData{
			Title: "Chat Lobby",
		},
		User: &pages.UserData{
			Username: user.Username,
		},
		Rooms: roomsList,
	}
	return c.Render(200, page.Template(), page)
}
