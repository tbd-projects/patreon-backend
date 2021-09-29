package general_handlers

import (
	"github.com/gorilla/mux"
	"patreon/internal/app"
)

const (
	GET      = "GET"
	POST     = "POST"
	OPTIONAL = "OPTIONAL"
	PUT      = "PUT"
)

type HandlerJoiner struct {
	joinedHandlers []app.Joinable
	currentUrl     string
}

func NewHandlerJoiner(joinedHandlers []app.Joinable, currentUrl string) *HandlerJoiner {
	return &HandlerJoiner{
		joinedHandlers: joinedHandlers,
		currentUrl:     currentUrl,
	}
}

func (h *HandlerJoiner) AddHandlers(joinedHandlers []app.Joinable) {
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
