package server

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/server/joinable_handler"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MainHandler struct {
	baseHandler joinable_handler.HandlerJoiner
	router      *mux.Router
	log         *logrus.Logger
}

func NewMainHandler(router *mux.Router, joinedHandlers []joinable_handler.IJoinable) *MainHandler {
	return &MainHandler{
		baseHandler: joinable_handler.CreateHandlerJoiner(joinedHandlers, ""),
		log:         logrus.New(),
		router:      router,
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

func (h *MainHandler) Join() {
	h.baseHandler.Join(h.router)
}
