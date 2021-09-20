package handlers

import (
	"net/http"

	"github.com/csothen/birdy"
	"github.com/gorilla/mux"
)

type Chat struct {
	service birdy.ChatService
}

func NewChat(s birdy.ChatService) *Chat {
	return &Chat{s}
}

func (c *Chat) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", c.getAll).Methods(http.MethodGet)
	r.HandleFunc("/", c.create).Methods(http.MethodGet)
	r.HandleFunc("/{id}", c.getOne).Methods(http.MethodGet)
	r.HandleFunc("/{id}", c.update).Methods(http.MethodPut)
	r.HandleFunc("/{id}", c.delete).Methods(http.MethodDelete)
}

func (c *Chat) getAll(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}

func (c *Chat) create(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}

func (c *Chat) getOne(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}

func (c *Chat) update(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}

func (c *Chat) delete(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}
