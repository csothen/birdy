package handlers

import (
	"net/http"

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

	u, token, err := h.authService.Authenticate(lr.Username, lr.Password)
	if err != nil {
		c.Logger().Warnf("could not authenticate user: %+v", err)
		return c.String(401, "invalid username and password")
	}

	c.SetCookie(&http.Cookie{
		Name:    authCookieName,
		Value:   token.Value,
		Expires: token.Expiration,
	})

	sRooms := h.chatService.ListRooms()
	rooms := []room{}
	for _, r := range sRooms {
		rooms = append(rooms, room{
			ID:   r.ID.String(),
			Name: r.Name,
		})
	}

	loggedInPage := indexPage{
		Metadata:   metadata{"Birdy | Lobby"},
		IsLoggedIn: true,
		User:       user{Username: u.Username},
		Rooms:      rooms,
	}
	return c.Render(200, "base", loggedInPage)
}
