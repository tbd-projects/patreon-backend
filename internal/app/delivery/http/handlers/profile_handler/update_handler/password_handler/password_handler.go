package password_handler

import (
	"encoding/json"
	"errors"
	"io"
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

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UpdatePasswordHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewUpdatePasswordHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *UpdatePasswordHandler {
	h := &UpdatePasswordHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMethod(http.MethodPut, h.PUT,
		csrf_middleware.NewCsrfMiddleware(log, usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken,
		middleware.NewSessionMiddleware(h.sessionManager, log).Check,
	)
	return h
}

// PUT ChangePassword
// @Summary set new user password
// @Description change password from user
// @Accept  json
// @Produce json
// @Param user body models.RequestChangePassword true "Request body for change password"
// @Success 200 "success update password"
// @Failure 400 {object} models.ErrResponse "incorrect new password"
// @Failure 403 "csrf token is invalid, get new token"
// @Failure 404 {object} models.ErrResponse "User not found"
// @Failure 418 "User are authorized"
// @Failure 422 {object} models.ErrResponse "Not valid body"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 500 {object} models.ErrResponse "database error"
// @Router /user/update/password [PUT]
func (h *UpdatePasswordHandler) PUT(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	req := &models.RequestChangePassword{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
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
	err := h.userUsecase.UpdatePassword(userId, req.NewPassword)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
