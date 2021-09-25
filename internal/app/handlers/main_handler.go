package handlers

import (
	"net/http"
	"patreon/internal/app"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type MainHandler struct {
	baseHandler app.HandlerJoiner
	router      *mux.Router
	RespondHandler
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		baseHandler:    app.HandlerJoiner{},
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.PrintRequest(r)
	h.router.ServeHTTP(w, r)
}
func (h *MainHandler) JoinHandlers(joinedHandlers []app.Joinable) {
	h.baseHandler.AddHandlers(joinedHandlers)
	h.baseHandler.Join(h.router)

}
func (h *MainHandler) SetRouter(router *mux.Router) {
	h.router = router
}
func (h *MainHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *MainHandler) Join(router *mux.Router) {
	h.baseHandler.Join(router)
}
