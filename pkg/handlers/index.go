package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Index(c echo.Context) error {
	notLoggedInPage := indexPage{
		Metadata:   metadata{"Birdy | Home"},
		IsLoggedIn: false,
	}

	cookie, err := c.Cookie(authCookieName)
	if err != nil {
		c.Logger().Errorf("could not retrieve cookie: %+v", err)
		return c.Render(200, "base", notLoggedInPage)
	}

	u, err := h.authService.Validate(cookie.Value)
	if err != nil {
		c.Logger().Errorf("error validating token: %+v", err)
		return c.Render(200, "base", notLoggedInPage)
	}

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
