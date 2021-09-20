package attachable_handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type IAttachable interface {
	Attach(router *mux.Router)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

type HandlerAttacher struct {
	attachedHandlers []IAttachable
	currentUrl       string
}

func (attacher HandlerAttacher) GetUrl() string {
	return attacher.currentUrl
}

func CreateHandlerAttacher(attachedHandlers []IAttachable, currentUrl string) HandlerAttacher {
	return HandlerAttacher{attachedHandlers, currentUrl}
}

func (attacher *HandlerAttacher) Attach(router *mux.Router) {
	subrouter := router.PathPrefix(attacher.currentUrl).Subrouter()

	for _, handler := range attacher.attachedHandlers {
		handler.Attach(subrouter)
	}
}
