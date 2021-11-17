package aw_id_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/middleware"
	sessionMid "patreon/internal/app/sessions/middleware"
	useAwards "patreon/internal/app/usecase/awards"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AwardsIdHandler struct {
	awardsUsecase useAwards.Usecase
	bh.BaseHandler
}

func NewAwardsIdHandler(log *logrus.Logger,
	ucAwards useAwards.Usecase, sClient session_client.AuthCheckerClient) *AwardsIdHandler {
	h := &AwardsIdHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		awardsUsecase: ucAwards,
	}

	h.AddMethod(http.MethodDelete, h.DELETE, sessionMid.NewSessionMiddleware(sClient, log).CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		middleware.NewAwardsMiddleware(log, ucAwards).CheckCorrectAwardFunc,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)

	return h
}

// DELETE Awards
// @Summary delete current awards
// @Description delete current awards from current creator
// @Produce json
// @Success 200
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 404 {object} models.ErrResponse "award with this id not found"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator", "this awards not belongs this creators", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
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
