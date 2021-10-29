package login_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/sessions/sessions_manager"
	usecase_user "patreon/internal/app/usecase/user"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewLoginHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase) *LoginHandler {
	h := &LoginHandler{
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
		sessionManager: sManager,
		userUsecase:    ucUser,
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).CheckNotAuthorized)
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST Login
// @Summary login user
// @Description login user
// @Accept  json
// @Produce json
// @Param user body models.RequestLogin true "Request body for login"
// @Success 200 "Successfully login"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 500 {object} models.ErrResponse "can not create session"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 401 {object} models.ErrResponse "incorrect email or password"
// @Failure 418 "User are authorized"
// @Router /login [POST]
func (h *LoginHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Fatal(err)
		}
	}(r.Body)

	req := &models.RequestLogin{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(req); err != nil ||
		len(req.Login) == 0 || len(req.Password) == 0 {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	h.Log(r).Debugf("Login : %s, password : %s", req.Login, req.Password)

	id, err := h.userUsecase.Check(req.Login, req.Password)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrors)
		return
	}

	res, err := h.sessionManager.Create(id)
	if err != nil || res.UserID != id {
		h.Log(r).Errorf("Error create session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateSession)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    res.UniqID,
		Expires:  time.Now().Add(sessions_manager.ExpiredCookiesTime),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
