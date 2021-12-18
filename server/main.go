package main

import (
	"fmt"
	"log"
	"os"

	"github.com/csothen/birdy/services/chat"
	"github.com/urfave/cli/v2"
)

var port string

type App struct {
	app *cli.App
}

func NewApp() *App {
	return &App{
		app: cli.NewApp(),
	}
}

func (a *App) Run() error {
	a.setInfo()
	a.setCommands()
	return a.app.Run(os.Args)
}

func main() {
	app := NewApp()
	log.Fatal(app.Run())
}

func (a *App) setInfo() {
	a.app.Name = "Birdy Chatting Server CLI"
	a.app.Usage = "CLI for starting the Birdy Chatting Server"
	a.app.Authors = []*cli.Author{
		{
			Name:  "César Pinheiro",
			Email: "cesarjcpinheiro@gmail.com",
		},
	}
	a.app.Version = "1.0.0"
}

func (a *App) setCommands() {
	a.app.Commands = []*cli.Command{
		{
			Name: "start",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "port",
					Aliases:     []string{"p"},
					Usage:       "Set the port in which the server will be served",
					Value:       "8080",
					Destination: &port,
				},
			},
			Usage: "Start the server",
			Action: func(c *cli.Context) error {
				srv := chat.NewServer()
				go srv.Run()

				addr := fmt.Sprintf(":%s", port)
				return srv.Start(addr)
			},
		},
	}
}
