package handlers

import "github.com/labstack/echo/v4"

type IndexPage struct {
	Head HeadData
}

type HeadData struct {
	Title string
}

func (h *Handler) Index(c echo.Context) error {
	return c.Render(200, "base", IndexPage{
		Head: HeadData{"Birdy"},
	})
}
