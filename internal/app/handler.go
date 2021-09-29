package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Error(w http.ResponseWriter, r *http.Request, code int, err error)
	Respond(w http.ResponseWriter, r *http.Request, code int, data interface{})
}

type Joinable interface {
	Join(router *mux.Router)
}

