package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	Store          store.Store
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/login"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *LoginHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *LoginHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *LoginHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *LoginHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST", "OPTIONS")
	h.baseHandler.Join(router)
}
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Fatal(err)
		}
	}(r.Body)

	req := &request{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil {
		h.log.Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	u, err := h.Store.User().FindByLogin(req.Login)
	h.log.Debugf("Login : %s, password : %s", req.Login, req.Password)
	if err != nil || !u.ComparePassword(req.Password) {
		h.log.Warnf("Fail get user or compare password %s", err)
		h.Error(w, r, http.StatusUnauthorized, handler_errors.IncorrectEmailOrPassword)
		return
	}

	res, err := h.SessionManager.Create(int64(u.ID))
	if err != nil || res.UserID != int64(u.ID) {
		h.log.Errorf("Error create session %s", err)
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
