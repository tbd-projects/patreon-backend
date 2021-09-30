package general_handlers

import (
	"github.com/gorilla/mux"
	"patreon/internal/app"
	"patreon/internal/app/handlers/urls"
)

const (
	GET      = "GET"
	POST     = "POST"
	OPTIONAL = "OPTIONAL"
	PUT      = "PUT"
)

type HandlerJoiner struct {
	joinedHandlers []app.Joinable
	handlerUrl     urls.UrlPath
}

func NewHandlerJoiner(joinedHandlers []app.Joinable, handlerUrl urls.UrlPath) *HandlerJoiner {
	return &HandlerJoiner{
		joinedHandlers: joinedHandlers,
		handlerUrl:     handlerUrl,
	}
}

func (h *HandlerJoiner) JoinHandlers(joinedHandlers ...app.Joinable) {
	for _, handler := range joinedHandlers {
		h.joinedHandlers = append(h.joinedHandlers, handler)
	}
}

func (h HandlerJoiner) GetUrl() urls.UrlPath {
	return h.handlerUrl
}

func (h *HandlerJoiner) Join(router *mux.Router) {
	subrouter := router.PathPrefix(string(h.handlerUrl)).Subrouter()

	for _, handler := range h.joinedHandlers {
		handler.Join(subrouter)
	}
}
