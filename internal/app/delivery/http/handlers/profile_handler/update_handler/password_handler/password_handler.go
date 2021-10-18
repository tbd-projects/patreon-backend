package password_handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"patreon/internal/app"
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
	h.AddMethod(http.MethodPut, h.PUT)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
	return h
}

// PUT ChangePassword
// @Summary set new user password
// @Description change password from user
// @Accept  json
// @Produce json
// @Param user body models.RequestChangePassword true "Request body for change password"
// @Success 200 {object} "success update password"
// @Failure 400 {object} "incorrect new password"
// @Failure 404 {object} "User not found"
// @Failure 418 "User are authorized"
// @Failure 422 {object} "Not valid body"
// @Failure 500 {object} "server error"
// @Failure 500 {object} "database error"
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
	h.Respond(w, r, http.StatusOK, map[string]string{"respond": "success update password"})
}
