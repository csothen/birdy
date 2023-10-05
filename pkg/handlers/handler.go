package handlers

import (
	"github.com/csothen/birdy/pkg/auth"
	"github.com/csothen/birdy/pkg/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	authCookieName = "auth-token"
	userContextKey = "user"
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

func (h *Handler) Protected(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(authCookieName)
		if err != nil {
			c.Logger().Errorf("could not retrieve cookie: %+v", err)
			return c.String(401, "unauthorized")
		}

		u, err := h.authService.Validate(cookie.Value)
		if err != nil {
			c.Logger().Errorf("error validating token: %+v", err)
			return c.String(401, "unauthorized")
		}

		c.Set(userContextKey, u)

		return next(c)
	}
}
