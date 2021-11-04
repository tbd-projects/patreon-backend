package subscribe_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	responseModels "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	middleSes "patreon/internal/app/sessions/middleware"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type SubscribeHandler struct {
	sessionManager    sessions.SessionsManager
	subscriberUsecase usecase_subscribers.Usecase
	bh.BaseHandler
}

func NewSubscribeHandler(log *logrus.Logger, sManager sessions.SessionsManager,
	ucSubscribers usecase_subscribers.Usecase) *SubscribeHandler {
	h := &SubscribeHandler{
		BaseHandler:       *bh.NewBaseHandler(log),
		subscriberUsecase: ucSubscribers,
		sessionManager:    sManager,
	}
	h.AddMethod(http.MethodGet, h.GET, middleSes.NewSessionMiddleware(h.sessionManager, log).CheckFunc)
	return h
}

// GET Subscribers
// @Summary subscribers of the creator
// @Description get subscribers of the creators with id = creator_id
// @Produce json
// @Param creator_id path int true "creator_id"
// @Success 200 {object} models.SubscribersCreatorResponse "Successfully get creator subscribers with creator id = creator_id"
// @Failure 400 {object} models.ErrResponse "invalid parameters - creator_id"
// @Failure 401 {object} models.ErrResponse "User are not authorized"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Router /creators/{:creator_id}/subscribers [GET]
func (h *SubscribeHandler) GET(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		h.Log(r).Warnf("invalid creator_id %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parameters %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	subscribers, err := h.subscriberUsecase.GetSubscribers(creatorID)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}
	res := responseModels.ToSubscribersCreatorResponse(subscribers)
	h.Log(r).Debugf("get users %v", subscribers)
	h.Respond(w, r, http.StatusOK, res)
}
