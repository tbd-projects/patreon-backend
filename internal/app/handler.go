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

type HandlerJoiner struct {
	joinedHandlers []Joinable
	currentUrl     string
}

func NewHandlerJoiner(joinedHandlers []Joinable, currentUrl string) *HandlerJoiner {
	return &HandlerJoiner{
		joinedHandlers: joinedHandlers,
		currentUrl:     currentUrl,
	}
}
func (h *HandlerJoiner) AddHandlers(joinedHandlers []Joinable) {
	for _, handler := range joinedHandlers {
		h.joinedHandlers = append(h.joinedHandlers, handler)
	}
}
func (h HandlerJoiner) GetUrl() string {
	return h.currentUrl
}

func (h *HandlerJoiner) Join(router *mux.Router) {
	subrouter := router.PathPrefix(h.currentUrl).Subrouter()

	for _, handler := range h.joinedHandlers {
		handler.Join(subrouter)
	}
}
