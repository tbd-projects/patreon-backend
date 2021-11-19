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
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type UpdateNicknameHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdateNicknameHandler(log *logrus.Logger,
	sClient session_client.AuthCheckerClient, ucUser usecase_user.Usecase) *UpdateNicknameHandler {
	h := &UpdateNicknameHandler{
		sessionClient: sClient,
		userUsecase:   ucUser,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMiddleware(session_middleware.NewSessionMiddleware(h.sessionClient, log).Check)

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
// @Failure 404 {object} models.ErrResponse "user not found"
// @Failure 409 {object} models.ErrResponse "nickname already exists"
// @Failure 422 {object} models.ErrResponse "invalid body in request| user with this oldNickname not found | invalid nickname in body | old nickname not equal current user nickname"
// @Failure 500 {object} models.ErrResponse "can not do bd operation", "server error"
// @Failure 418 "User are authorized"
// @Router /user/update/nickname [PUT]
func (h *UpdateNicknameHandler) PUT(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("user_id").(int64)
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
	err := h.userUsecase.UpdateNickname(id, req.OldNickname, req.NewNickname)
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}
	w.WriteHeader(http.StatusOK)
}
