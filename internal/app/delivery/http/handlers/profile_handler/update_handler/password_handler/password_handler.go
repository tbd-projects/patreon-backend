package password_handler

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

type UpdatePasswordHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdatePasswordHandler(log *logrus.Logger,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *UpdatePasswordHandler {
	h := &UpdatePasswordHandler{
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

// PUT ChangePassword
// @Summary set new user password
// @Description change password from user
// @Accept  json
// @Produce json
// @Param password body http_models.RequestChangePassword true "Request body for change password"
// @Success 200 "success update password"
// @Failure 409 {object} http_models.ErrResponse "incorrect new password"(mean old password equal new)
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token", "incorrect email or password"
// @Failure 404 {object} http_models.ErrResponse "user not found"
// @Failure 418 "User are authorized"
// @Failure 400 {object} http_models.ErrResponse "invalid body in request", "incorrect new password"
// @Failure 500 {object} http_models.ErrResponse "server error", "can not do bd operation"
// @Router /user/update/password [PUT]
func (h *UpdatePasswordHandler) PUT(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestChangePassword{}
	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil || req.OldPassword == "" || req.NewPassword == "" {
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidBody)
		return
	}
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.HandlerError(w, r, http.StatusInternalServerError, app.GeneralError{
			Err:         handler_errors.InternalError,
			ExternalErr: errors.New("context parse userId error"),
		})
		return
	}
	err = h.userUsecase.UpdatePassword(userId, req.OldPassword, req.NewPassword)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
