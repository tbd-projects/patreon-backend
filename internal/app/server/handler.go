package server

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/store"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
	router *Router
	store  *store.Store
	log    *logrus.Logger
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		log: logrus.New(),
	}
}

func (h *MainHandler) SetRouter(router *Router) {
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
