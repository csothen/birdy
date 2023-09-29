package main

import (
	"html/template"
	"io"
	"log/slog"
	"os"

	"github.com/csothen/birdy/pkg/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	tmpls, err := template.ParseGlob("public/views/*.html")
	if err != nil {
		logger.Error("couldn't initialize templates: %v", err)
		return
	}

	handler := handlers.NewHandler()

	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())

	e.GET("/", handler.Index)
	e.GET("/join/:room", handler.JoinRoom)

	e.Logger.Fatal(e.Start(":8080"))
}
