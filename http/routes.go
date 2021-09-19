package http

import (
	"net/http"

	"github.com/csothen/birdy/handlers"
	"github.com/gorilla/mux"
)

func (a *Application) registerRoutes(r *mux.Router) {
	r.HandleFunc("/auth/login", handlers.HandleLogin).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", handlers.HandleLogout).Methods(http.MethodPost)
	r.HandleFunc("/chats", handlers.HandleGetAllChats).Methods(http.MethodGet)
}
