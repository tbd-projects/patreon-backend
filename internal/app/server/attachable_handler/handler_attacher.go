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
	handlerFunc      HandlerFunc
}

func createHandlerAttacher(attachedHandlers []IAttachable, currentUrl string,
	handlerFunc HandlerFunc) *HandlerAttacher {
	return &HandlerAttacher{attachedHandlers, currentUrl, handlerFunc}
}

func (attacher *HandlerAttacher) Attach(router *mux.Router) {
	router.HandleFunc(attacher.currentUrl + "/", attacher.handlerFunc)
	subrouter := router.PathPrefix(attacher.currentUrl + "/").Subrouter()

	for _, handler := range attacher.attachedHandlers {
		handler.Attach(subrouter)
	}
}
