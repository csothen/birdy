package handlers

import (
	"github.com/csothen/birdy/pkg/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Index(c echo.Context) error {
	page := pages.IndexPage{
		Head: pages.HeadData{Title: "Birdy"},
	}

	return c.Render(200, page.Template(), page)
}
