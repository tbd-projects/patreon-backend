package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/server/attachable_handler"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
	baseHandler attachable_handler.HandlerAttacher
	router      *mux.Router
	log         *logrus.Logger
}

func NewMainHandler(router *mux.Router,	attachedHandlers []attachable_handler.IAttachable) *MainHandler {
	return &MainHandler{
		baseHandler: attachable_handler.CreateHandlerAttacher(attachedHandlers, ""),
		log:    logrus.New(),
		router: router,
	}
}

func (h *MainHandler) SetRouter(router *mux.Router) {
	h.router = router
}

func (h *MainHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *MainHandler) Attach() {
	h.baseHandler.Attach(h.router)
}
