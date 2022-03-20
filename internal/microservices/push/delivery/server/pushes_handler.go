package push_server

import (
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"
	"patreon/internal/microservices/push/push/usecase"
)

type PushesHandler struct {
	sessionClient session_client.AuthCheckerClient
	usecase       usecase.Usecase
	bh.BaseHandler
}

func NewPushesHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient, usecase usecase.Usecase) *PushesHandler {
	h := &PushesHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sManager,
		usecase:       usecase,
	}
	h.AddMethod(http.MethodGet, h.GET, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET user pushes
// @Summary get user pushes
// @Description get user pushes
// @Produce json
// @tags user
// @Success 200 {object} PushesResponse
// @Failure 500 {object} http_models.ErrResponse "server error"
// @Failure 401 "user are not authorized"
// @Router /user/pushes [GET]
func (h *PushesHandler) GET(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Errorf("not found user_id in context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	res, err := h.usecase.GetPushInfo(userId)
	if err != nil {
		h.Log(r).Errorf("server error: %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	h.Respond(w, r, http.StatusOK, ToPushesResponse(res))
}
