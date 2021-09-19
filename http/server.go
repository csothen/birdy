package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/csothen/birdy"
	"github.com/csothen/birdy/data"
	"github.com/csothen/env"
	"github.com/gorilla/mux"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second

	defaultAddress = ":8080"
)

type Application struct {
	server *http.Server
	router *mux.Router

	// Services used by the routes
	AuthService birdy.AuthService
	UserService birdy.UserService
	ChatService birdy.ChatService
}

func NewServer() *Application {
	l := log.New(os.Stdout, "[ birdy ] ", log.LstdFlags)
	p := env.NewParser(l)

	app := &Application{
		server: &http.Server{
			Addr:         p.String("BIND_ADDRESS", defaultAddress),
			ErrorLog:     l,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		router: mux.NewRouter(),
	}

	// Middlewares
	app.router.Use(handlePanicMiddleware)
	app.router.Use(jsonMiddleware)

	app.server.Handler = http.HandlerFunc(app.serveHTTP)
	app.router.NotFoundHandler = http.HandlerFunc(app.handleNotFound)

	sr := app.router.PathPrefix("/").Subrouter()
	app.registerRoutes(sr)

	return app
}

func (a *Application) serveHTTP(rw http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(rw, r)
}

func (*Application) handleNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	data.ToJSON(&birdy.GenericError{
		Message: fmt.Sprintf("The route '%s' does not exist", r.URL.Path),
	}, rw)
	return
}
