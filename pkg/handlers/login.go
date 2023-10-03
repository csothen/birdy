package handlers

import (
	"net/http"
	"time"

	"github.com/csothen/birdy/pkg/pages"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c echo.Context) error {
	var lr LoginRequest

	if err := c.Bind(&lr); err != nil {
		page := pages.ErrorPage{
			Code:  http.StatusBadRequest,
			Error: err,
		}

		return c.Render(200, page.Template(), page)
	}

	c.SetCookie(&http.Cookie{
		Name:    "authentication-token",
		Value:   "token",
		Expires: time.Now().Add(30 * time.Minute),
	})

	page := pages.IndexPage{
		Head: pages.HeadData{
			Title: "Chat Lobby",
		},
		User: &pages.UserData{
			Username: "username",
		},
		Rooms: []pages.RoomData{},
	}
	return c.Render(200, page.Template(), page)
}
