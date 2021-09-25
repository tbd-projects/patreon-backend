package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
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
	router         *mux.Router
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
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST", "GET")

	//router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	//router.Use(h.authMiddleware.CheckNotAuthorized)
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
		h.Error(h.log, w, r, http.StatusUnprocessableEntity, err)
		return
	}
	u, err := h.Store.User().FindByLogin(req.Login)

	if err != nil || !u.ComparePassword(req.Password) {
		h.Error(h.log, w, r, http.StatusUnauthorized, store.IncorrectEmailOrPassword)
		return
	}
	res, err := h.SessionManager.Create(int64(u.ID))
	if err != nil || res.UserID != int64(u.ID) {
		h.Error(h.log, w, r, http.StatusInternalServerError, store.IncorrectEmailOrPassword)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    res.UniqID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	h.Respond(h.log, w, r, http.StatusOK, "successfully login")
}
