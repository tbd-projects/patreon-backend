package login_handler

import (
	"encoding/json"
	"io"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"
	"time"

	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewLoginHandler(log *logrus.Logger, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase) *LoginHandler {
	h := &LoginHandler{
		BaseHandler:    *bh.NewBaseHandler(log),
		sessionManager: sManager,
		userUsecase:    ucUser,
	}
	h.AddMethod(http.MethodPost, h.POST)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, h.Log()).CheckNotAuthorized)
	return h
}

// Login
// @Summary login user
// @Description login user
// @Accept  json
// @Produce json
// @Param user body models.RequestLogin true "Request body for login"
// @Success 201 {object} models.BaseResponse "Successfully login"
// @Failure 401 {object} models.BaseResponse "Incorrect password or email"
// @Failure 422 {object} models.BaseResponse "Not valid body"
// @Failure 500 {object} models.BaseResponse "Creation error in sessions"
// @Failure 418 "User are authorized"
// @Router /login [POST]
func (h *LoginHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Fatal(err)
		}
	}(r.Body)

	req := &models.RequestLogin{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil ||
		len(req.Login) == 0 || len(req.Password) == 0 {
		h.Log().Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	h.Log().Debugf("Login : %s, password : %s", req.Login, req.Password)
	id, err := h.userUsecase.Check(req.Login, req.Password)
	if err != nil {
		h.Log().Warnf("Fail check user %s", err)
		h.UsecaseError(w, r, err, codesByErrors)
		//h.Error(w, r, http.StatusUnauthorized, handler_errors.IncorrectEmailOrPassword)
		return
	}

	res, err := h.sessionManager.Create(id)
	if err != nil || res.UserID != id {
		h.Log().Errorf("Error create session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateSession)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    res.UniqID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	h.Respond(w, r, http.StatusOK, "successfully login")
}
