package http

import (
	"log"
	"net/http"
	"time"

	"github.com/csothen/birdy"
	"github.com/csothen/birdy/data"
	"github.com/csothen/birdy/handlers"
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

	// Handlers
	AuthHandler birdy.AuthHandler
	ChatHandler birdy.ChatHandler
}

func NewServer(l *log.Logger, as birdy.AuthService, cs birdy.ChatService) *Application {
	p := env.NewParser(l)

	// Create handlers
	ah := handlers.NewAuth(as)
	ch := handlers.NewChat(cs)

	app := &Application{
		server: &http.Server{
			Addr:         p.String("BIND_ADDRESS", defaultAddress),
			ErrorLog:     l,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		router:      mux.NewRouter(),
		AuthHandler: ah,
		ChatHandler: ch,
	}

	// Middlewares
	app.router.Use(handlePanicMiddleware)
	app.router.Use(jsonMiddleware)

	app.server.Handler = http.HandlerFunc(app.serveHTTP)
	app.router.NotFoundHandler = http.HandlerFunc(app.handleNotFound)

	// Register routes
	sr := app.router.PathPrefix("/api").Subrouter()
	app.AuthHandler.RegisterRoutes(sr.PathPrefix("/auth").Subrouter())
	app.ChatHandler.RegisterRoutes(sr.PathPrefix("/chats").Subrouter())

	return app
}

func (a *Application) serveHTTP(rw http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(rw, r)
}

func (*Application) handleNotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	data.ToJSON(&birdy.GenericError{Message: "404 Not Found"}, rw)
}
