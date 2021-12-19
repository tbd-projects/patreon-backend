package push_server

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"
	"patreon/internal/microservices/push/push/repository"
	"patreon/internal/microservices/push/push/usecase"
)

type MarkPushHandler struct {
	sessionClient session_client.AuthCheckerClient
	usecase       usecase.Usecase
	bh.BaseHandler
}

func NewMarkPushHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient, usecase usecase.Usecase) *MarkPushHandler {
	h := &MarkPushHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sManager,
		usecase:       usecase,
	}
	h.AddMethod(http.MethodPut, h.PUT, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// PUT mark push as viewed
// @Summary mark push as viewed
// @Description mark push as viewed
// @Produce json
// @tags user
// @Success 200 {object} utils.PushResponse Type can be "Comment", "Post", "Payment"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 404 {object} http_models.ErrResponse "content not modify"
// @Failure 401 "user are not authorized"
// @Router /user/push/{:push_id} [PUT]
func (h *MarkPushHandler) PUT(w http.ResponseWriter, r *http.Request) {
	pushId, ok := h.GetInt64FromParam(w, r, "push_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Errorf("not found user_id in context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	err := h.usecase.MarkViewed(pushId, userId)
	if err != nil {
		h.Log(r).Errorf("server error: %s", err)
		if errors.Is(err, repository.NotModify) {
			h.Error(w, r, http.StatusNotFound, handler_errors.NoModify)
			return
		}
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
