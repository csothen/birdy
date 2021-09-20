package handlers

import (
	"net/http"

	"github.com/csothen/birdy"
	"github.com/gorilla/mux"
)

type Auth struct {
	service birdy.AuthService
}

func NewAuth(s birdy.AuthService) *Auth {
	return &Auth{s}
}

func (a *Auth) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", a.login).Methods(http.MethodPost)
	r.HandleFunc("/logout", a.logout).Methods(http.MethodPost)
}

func (a *Auth) login(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}

func (a *Auth) logout(rw http.ResponseWriter, r *http.Request) {
	// TODO: Handle request and call the service
}
