package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Connect(router *mux.Route)
}
