package subscribe_handler

import (
	"io"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
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

func NewSubscribeHandler(log *logrus.Logger, router *mux.Router,
	cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucSubscribers usecase_subscribers.Usecase) *SubscribeHandler {
	h := &SubscribeHandler{
		BaseHandler:       *bh.NewBaseHandler(log, router, cors),
		subscriberUsecase: ucSubscribers,
		sessionManager:    sManager,
	}
	h.AddMethod(http.MethodPost, h.POST, middleSes.NewSessionMiddleware(h.sessionManager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	h.AddMethod(http.MethodPost, h.DELETE, middleSes.NewSessionMiddleware(h.sessionManager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
	)
	return h
}

// DELETE Unsubscribe
// @Summary unsubscribe from the creator
// @Description unsubscribe from the creator with id = creator_id
// @Produce json
// @Success 200 "Successfully unsubscribe on the creator with id = creator_id"
// @Failure 400 "invalid parameters - creator_id"
// @Failure 401 "User are not authorized"
// @Failure 409 "this user is not subscribed on the creator"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Router /creators/{:creator_id}/unsubscribe [DELETE]
func (h *SubscribeHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	creatorID, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	subscriber := &models.Subscriber{
		UserID:    userID.(int64),
		CreatorID: creatorID,
	}
	err := h.subscriberUsecase.UnSubscribe(subscriber)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	h.Log(r).Debugf("unsubscribe from creator_id = %v", creatorID)
	w.WriteHeader(http.StatusOK)
}

// POST Subscribe
// @Summary subscribes on the creator
// @Description subscribes on the creator with id = creator_id
// @Produce json
// @Success 201 "Successfully subscribe on the creator with id = creator_id"
// @Failure 400 "invalid parameters - creator_id"
// @Failure 401 "User are not authorized"
// @Failure 409 "this user already subscribed on the creator"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Router /creators/{:creator_id}/subscribe [POST]
func (h *SubscribeHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	creatorID, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}
	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	subscriber := &models.Subscriber{
		UserID:    userID.(int64),
		CreatorID: creatorID,
	}
	err := h.subscriberUsecase.Subscribe(subscriber)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	h.Log(r).Debugf("subscribe on creator_id = %v", creatorID)
	w.WriteHeader(http.StatusCreated)
}
