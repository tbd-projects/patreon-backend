package nickname_handler

import (
	"errors"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type UpdateNicknameHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdateNicknameHandler(log *logrus.Logger,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *UpdateNicknameHandler {
	h := &UpdateNicknameHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log),
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)

	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)
	return h
}

// PUT NicknameChange
// @Summary set new user nickname
// @Accept  json
// @Param nickname body models.RequestChangeNickname true "Request body for change nickname"
// @Success 200 "successfully change nickname"
// @Failure 403 {object} models.ErrResponse "csrf token is invalid, get new token"
// @Failure 409 {object} models.ErrResponse "nickname already exists"
// @Failure 422 {object} models.ErrResponse "invalid body in request", "user with this oldNickname not found", "invalid nickname in body"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 418 "User are authorized"
// @Router /user/update/avatar [PUT]
func (h *UpdateNicknameHandler) PUT(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: errors.New("UpdateNicknameHandler: context parse userId error"),
		})
		return
	}
	req := &models.RequestChangeNickname{}
	if err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy()); err != nil {
		h.Log(r).Warnf("UpdateNicknameHandler: can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	if err := req.Validate(); err != nil {
		h.Log(r).Warnf("UpdateNicknameHandler: invalid body on request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, models.NicknameValidateError)
		return
	}
	err := h.userUsecase.UpdateNickname(req.OldNickname, req.NewNickname)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}
	w.WriteHeader(http.StatusOK)
}
