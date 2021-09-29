package handlers

import (
	"net/http"
	"patreon/internal/app"

	gh "patreon/internal/app/handlers/general_handlers"

	"github.com/gorilla/mux"
)

type MainHandler struct {
	joinedHandler gh.HandlerJoiner
	router        *mux.Router
	gh.RespondHandler
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		joinedHandler: gh.HandlerJoiner{},
	}
}

func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *MainHandler) JoinHandlers(joinedHandlers []app.Joinable) {
	h.joinedHandler.AddHandlers(joinedHandlers)
	h.joinedHandler.Join(h.router)
}

func (h *MainHandler) SetRouter(router *mux.Router) {
	h.router = router
}

func (h *MainHandler) Join(router *mux.Router) {
	h.joinedHandler.Join(router)
}
