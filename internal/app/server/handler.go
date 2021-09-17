package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app/store"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
	router *mux.Router
	store  *store.Store
	log    *logrus.Logger
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		log: logrus.New(),
	}
}

func (h *MainHandler) SetRouter(router *mux.Router) {
	h.router = router
}
func (h *MainHandler) SetStore(store *store.Store) {
	h.store = store
}
func (h *MainHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *MainHandler) HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "hello patron!")
		if err != nil {
			h.log.Fatal(err)
		}
	}
}

func (h *MainHandler) RegisterHandlers() {
	h.router.HandleFunc("/hello", h.HandleRoot()).Methods("GET", "POST")

}
