package aw_id_handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"
)

type AwardsIdHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsIdHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucAwards useAwards.Usecase, manager sessions.SessionsManager) *AwardsIdHandler {
	h := &AwardsIdHandler{
		BaseHandler:   *bh.NewBaseHandler(log, router, cors),
		awardsUsecase: ucAwards,
	}
	h.AddMethod(http.MethodDelete, h.DELETE, sessionMid.NewSessionMiddleware(manager, log).CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc)
	return h
}

// DELETE Awards
// @Summary delete current awards
// @Description delete current awards from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "award with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Failure 403 {object} models.ErrResponse "this awards not belongs this creators"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/awards/{:award_id} [DELETE]
func (h *AwardsIdHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	awardsId, ok := h.GetInt64FromParam(w, r, "award_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 2 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	err := h.awardsUsecase.Delete(awardsId)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsDELETE)
		return
	}

	w.WriteHeader(http.StatusOK)
}
