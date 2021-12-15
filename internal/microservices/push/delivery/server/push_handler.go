package push_server

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"
	"patreon/internal/microservices/push/utils"
)


type PushHandler struct {
	sessionClient session_client.AuthCheckerClient
	hub           *utils.SendHub
	upgrader      *websocket.Upgrader
	bh.BaseHandler
}

func NewPushHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient, hub *utils.SendHub,
	upgrader *websocket.Upgrader) *PushHandler {
	h := &PushHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sManager,
		hub:           hub,
		upgrader: upgrader,
	}
	h.AddMethod(http.MethodGet, h.GET, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET Push Creator
// @Summary create creator
// @Description create creator with id from path, and respond created creator
// @Param creator body http_models.RequestCreator true "Request body for creators"
// @Produce json
// @tags creators
// @Success 201 {object} http_models.IdResponse
// @Failure 409 {object} http_models.ErrResponse "creator already exist"
// @Failure 404 {object} http_models.ErrResponse "user with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid creator nickname", "invalid creator category-description", "invalid creator category", "invalid body in request"
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /user/push [GET]
func (h *PushHandler) GET(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Errorf("not found user_id in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := utils.NewClient(h.hub, userId, conn, h.Log(r))
	h.hub.RegisterClient(client)
	go client.SenderProcesses()
}
