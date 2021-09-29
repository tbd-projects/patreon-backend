package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"
	"patreon/internal/models"
	"time"

	gh "patreon/internal/app/handlers/general_handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	authMiddleware middleware.SessionMiddleware
	Store          store.Store
	SessionManager sessions.SessionsManager
	gh.RespondHandler
	withHideMethod
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		RespondHandler: gh.RespondHandler{},
		withHideMethod:    withHideMethod{gh.NewBaseHandler(logrus.New(), "/login")},
	}
}

func (h *LoginHandler) SetStore(store store.Store) {
	h.Store = store
}

func (h *LoginHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.Log())
}

func (h *LoginHandler) Join(router *mux.Router) {
	h.baseHandler.AddMethod(gh.POST, h.ServeHTTP)
	h.baseHandler.AddMethod(gh.GET, h.ServeHTTP)
	h.baseHandler.AddMiddleware(h.authMiddleware.CheckNotAuthorized)
	h.baseHandler.Join(router)
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
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Fatal(err)
		}
	}(r.Body)

	req := &models.RequestLogin{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil {
		h.Log().Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	u, err := h.Store.User().FindByLogin(req.Login)
	h.Log().Debugf("Login : %s, password : %s", req.Login, req.Password)
	if err != nil || !u.ComparePassword(req.Password) {
		h.Log().Warnf("Fail get user or compare password %s", err)
		h.Error(w, r, http.StatusUnauthorized, handler_errors.IncorrectEmailOrPassword)
		return
	}

	res, err := h.SessionManager.Create(int64(u.ID))
	if err != nil || res.UserID != int64(u.ID) {
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
