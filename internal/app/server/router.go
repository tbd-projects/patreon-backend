package server

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

func (mr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mr.router.ServeHTTP(w, r)
}

func (r *Router) HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello patron!")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r *Router) Configure() {
	r.router.HandleFunc("/hello", r.HandleRoot()).Methods("GET", "POST")
}
