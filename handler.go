package birdy

import "github.com/gorilla/mux"

type Handler interface {
	RegisterRoutes(router *mux.Router)
}

type AuthHandler interface {
	Handler
}

type ChatHandler interface {
	Handler
}
