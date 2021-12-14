package aw_subscribe_handler

import (
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	"patreon/internal/app/models"
	useAwards "patreon/internal/app/usecase/awards"
	usecase_pay_token "patreon/internal/app/usecase/pay_token"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AwardsSubscribeHandler struct {
	subscriberUsecase usecase_subscribers.Usecase
	awardsUsecase     useAwards.Usecase
	payTokenUsecase   usecase_pay_token.Usecase
	bh.BaseHandler
}

func NewAwardsSubscribeHandler(log *logrus.Logger, sClient session_client.AuthCheckerClient,
	ucSubscribers usecase_subscribers.Usecase, ucAwards useAwards.Usecase) *AwardsSubscribeHandler {
	h := &AwardsSubscribeHandler{
		BaseHandler:       *bh.NewBaseHandler(log),
		subscriberUsecase: ucSubscribers,
	}
	h.AddMethod(http.MethodPost, h.POST, session_middleware.NewSessionMiddleware(sClient, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
	)
	h.AddMethod(http.MethodDelete, h.DELETE, session_middleware.NewSessionMiddleware(sClient, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
	)

	return h
}

// DELETE Unsubscribe
// @Summary unsubscribe from the creator
// @tags awards
// @Description unsubscribe from the creator with id = creator_id and awards_id = award_id
// @Produce json
// @Param award_id path int true "award_id"
// @Param creator_id path int true "creator_id"
// @Success 200 "Successfully unsubscribe on the creator with id = creator_id"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters"
// @Failure 404 {object} http_models.ErrResponse "award with this id not found"
// @Failure 409 {object} http_models.ErrResponse "subscribes on the creator not found"
// @Failure 500 {object} http_models.ErrResponse "server error", "can not do bd operation"
// @Failure 403 {object} http_models.ErrResponse "this awards not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
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
// @tags awards
// @Description subscribes on the creator with id = creator_id
// @Accept json
// @Produce json
// @Param award_id path int true "award_id"
// @Param creator_id path int true "creator_id"
// @Param pay_token body http_models.SubscribeRequest true "Request payToken"
// @Success 201 "Successfully subscribe on the creator with id = creator_id"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "this user was not given this token"
// @Failure 409 {object} http_models.ErrResponse "this user already have subscribe on creator"
// @Failure 404 {object} http_models.ErrResponse "award with this id not found"
// @Failure 500 {object} http_models.ErrResponse "server error", "can not do bd operation"
// @Failure 403 {object} http_models.ErrResponse "this awards not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/awards/{:award_id}/subscribe [POST]
func (h *AwardsSubscribeHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.SubscribeRequest{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil || req.Validate() != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
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

	err = h.payTokenUsecase.CheckTokenByUser(models.PayToken{Token: req.Token}, userID.(int64))
	if err != nil {
		h.Log(r).Warnf("checkTokenVyUser error; err = %s, token from req = %s, userID = %v", err, req.Token, userID)
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	err = h.subscriberUsecase.Subscribe(subscriber, req.Token)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}
	h.Log(r).Debugf("subscribe on creator_id = %v", creatorID)
	w.WriteHeader(http.StatusCreated)
}
