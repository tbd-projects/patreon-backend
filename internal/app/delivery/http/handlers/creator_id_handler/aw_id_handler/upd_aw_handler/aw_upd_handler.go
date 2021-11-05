package aw_upd_handler

import (
	"image/color"
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	bd_modle "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

type AwardsUpdHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsUpdHandler(log *logrus.Logger,
	ucAwards useAwards.Usecase, manager sessions.SessionsManager) *AwardsUpdHandler {
	h := &AwardsUpdHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		awardsUsecase: ucAwards,
	}

	h.AddMethod(http.MethodPut, h.PUT, sessionMid.NewSessionMiddleware(manager, log).CheckFunc,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
	)

	return h
}

// PUT Awards
// @Summary update current awards
// @Description update current awards from current creator
// @Param award body models.RequestAwards true "Request body for update awards"
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "award with this id not found"
// @Failure 422 {object} models.ErrResponse "invalid body in request", "incorrect value of price", "empty name in request"
// @Failure 409 {object} models.ErrResponse "awards with this name already exists"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this awards not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/awards/{:award_id}/update [PUT]
func (h *AwardsUpdHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &models.RequestAwards{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	awardsId, ok := h.GetInt64FromParam(w, r, "award_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}
	award := &bd_modle.Award{
		ID:          awardsId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Color:       color.RGBA{R: req.Color.R, B: req.Color.B, G: req.Color.G, A: req.Color.A},
	}

	err = h.awardsUsecase.Update(award)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
