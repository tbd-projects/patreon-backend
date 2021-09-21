package joinable_handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type IJoinable interface {
	Join(router *mux.Router)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

type HandlerJoiner struct {
	joinedHandlers []IJoinable
	currentUrl     string
}

func (joiner HandlerJoiner) GetUrl() string {
	return joiner.currentUrl
}

func CreateHandlerJoiner(joinedHandlers []IJoinable, currentUrl string) HandlerJoiner {
	return HandlerJoiner{joinedHandlers, currentUrl}
}

func (joiner *HandlerJoiner) Join(router *mux.Router) {
	subrouter := router.PathPrefix(joiner.currentUrl).Subrouter()

	for _, handler := range joiner.joinedHandlers {
		handler.Join(subrouter)
	}
}
