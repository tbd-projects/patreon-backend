package aw_upd_handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"image/color"
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	bd_modle "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"
)

type AwardsUpdHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsUpdHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucAwards useAwards.Usecase, manager sessions.SessionsManager) *AwardsUpdHandler {
	h := &AwardsUpdHandler{
		BaseHandler:   *bh.NewBaseHandler(log, router, cors),
		awardsUsecase: ucAwards,
	}
	h.AddMethod(http.MethodPut, h.PUT, sessionMid.NewSessionMiddleware(manager, log).CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc)
	return h
}

// PUT Awards
// @Summary update current awards
// @Description update current awards from current creator
// @Param user body models.RequestAwards true "Request body for update awards"
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "award with this id not found"
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Failure 403 {object} models.ErrResponse "this post not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Failure 422 {object} models.ErrResponse "empty name in request"
// @Failure 422 {object} models.ErrResponse "incorrect value of price"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/awards/{:award_id}/update/other [PUT]
func (h *AwardsUpdHandler) PUT(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	awardsId, ok := h.GetInt64FromParam(w, r, "award_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	req := &models.RequestAwards{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	award := &bd_modle.Award{
		ID:          awardsId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Color:       color.RGBA{R: req.Color.R, B: req.Color.B, G: req.Color.G, A: req.Color.A},
	}

	err := h.awardsUsecase.Update(award)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
