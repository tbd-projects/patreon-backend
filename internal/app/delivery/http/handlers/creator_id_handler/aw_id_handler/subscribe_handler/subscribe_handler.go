package aw_subscribe_handler

import (
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	"patreon/internal/app/sessions"
	middleSes "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AwardsSubscribeHandler struct {
	subscriberUsecase usecase_subscribers.Usecase
	awardsUsecase     useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsSubscribeHandler(log *logrus.Logger, router *mux.Router,
	cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucSubscribers usecase_subscribers.Usecase, ucAwards useAwards.Usecase) *AwardsSubscribeHandler {
	h := &AwardsSubscribeHandler{
		BaseHandler:       *bh.NewBaseHandler(log, router, cors),
		subscriberUsecase: ucSubscribers,
	}
	h.AddMethod(http.MethodPost, h.POST, middleSes.NewSessionMiddleware(sManager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
	)
	h.AddMethod(http.MethodDelete, h.DELETE, middleSes.NewSessionMiddleware(sManager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
	)

	return h
}

// DELETE Unsubscribe
// @Summary unsubscribe from the creator
// @Description unsubscribe from the creator with id = creator_id and awards_id = award_id
// @Produce json
// @Param award_id path int true "award_id"
// @Param creator_id path int true "creator_id"
// @Success 200 "Successfully unsubscribe on the creator with id = creator_id"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 403 {object} models.ErrResponse "invalid csrf token"
// @Failure 403 {object} models.ErrResponse "incorrect award for creator"
// @Failure 404 {object} models.ErrResponse "award with this id not found"
// @Failure 401 {object} models.ErrResponse "User are not authorized"
// @Failure 409 {object} models.ErrResponse "this user is not subscribed on the creator"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Router /creators/{:creator_id}/awards/{:award_id}/subscribe [DELETE]
func (h *AwardsSubscribeHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	creatorID, _ := h.GetInt64FromParam(w, r, "creator_id")
	awardID, _ := h.GetInt64FromParam(w, r, "award_id")
	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	subscriber := &models.Subscriber{
		UserID:    userID.(int64),
		CreatorID: creatorID,
		AwardID:   awardID,
	}
	err := h.subscriberUsecase.UnSubscribe(subscriber)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}
	h.Log(r).Debugf("unsubscribe from creator_id = %v", creatorID)
	w.WriteHeader(http.StatusOK)
}

// POST Subscribe
// @Summary subscribes on the creator
// @Description subscribes on the creator with id = creator_id
// @Accept json
// @Produce json
// @Param award_id path int true "award_id"
// @Param creator_id path int true "creator_id"
// @Success 201 "Successfully subscribe on the creator with id = creator_id"
// @Failure 400 {object} models.ErrResponse "invalid parameters - creator_id"
// @Failure 401 {object} models.ErrResponse "User are not authorized"
// @Failure 403 {object} models.ErrResponse "invalid csrf token"
// @Failure 403 {object} models.ErrResponse "incorrect award for creator"
// @Failure 409 {object} models.ErrResponse "this user already subscribed on the creator"
// @Failure 500 {object} models.ErrResponse "serverError"
// @Router /creators/{:creator_id}/awards/{:award_id}/subscribe [POST]
func (h *AwardsSubscribeHandler) POST(w http.ResponseWriter, r *http.Request) {
	//req := &responseModels.SubscribeRequest{}
	//
	//err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	//if err != nil || req.Validate() != nil {
	//	h.Log(r).Warnf("can not parse request %s", err)
	//	h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
	//	return
	//}
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	creatorID, _ := h.GetInt64FromParam(w, r, "creator_id")
	awardID, _ := h.GetInt64FromParam(w, r, "award_id")

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	subscriber := &models.Subscriber{
		UserID:    userID.(int64),
		CreatorID: creatorID,
		AwardID:   awardID,
	}

	err := h.subscriberUsecase.Subscribe(subscriber)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	h.Log(r).Debugf("subscribe on creator_id = %v", creatorID)
	w.WriteHeader(http.StatusCreated)
}
