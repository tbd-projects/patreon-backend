package subscriptions_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	responseModels "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"

	"github.com/sirupsen/logrus"
)

type SubscriptionsHandler struct {
	sessionManager     sessions.SessionsManager
	subscribersUsecase usecase_subscribers.Usecase
	bh.BaseHandler
}

func NewSubscriptionsHandler(log *logrus.Logger, sManager sessions.SessionsManager,
	ucSubscribers usecase_subscribers.Usecase) *SubscriptionsHandler {
	h := &SubscriptionsHandler{
		sessionManager:     sManager,
		subscribersUsecase: ucSubscribers,
		BaseHandler:        *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		middleware.NewSessionMiddleware(h.sessionManager, log).CheckFunc)
	return h
}

// GET Subscriptions
// @Summary get user subscriptions
// @Description get user creators
// @Produce json
// @Success 200 {object} models.SubscriptionsUserResponse "Successfully get user subscriptions"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Failure 401 "user are not authorized"
// @Router /user/subscriptions [GET]
func (h SubscriptionsHandler) GET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	creators, err := h.subscribersUsecase.GetCreators(userID.(int64))
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	res := responseModels.ToSubscriptionsUser(creators)
	h.Log(r).Debugf("get creators %v", creators)
	h.Respond(w, r, http.StatusOK, res)
}
