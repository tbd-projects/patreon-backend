package handlers

import (
	"errors"
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

type LogoutHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	router         *mux.Router
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/logout"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *LogoutHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *LogoutHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *LogoutHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("GET")
	h.baseHandler.Join(router)
}
func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Fatal(err)
		}
	}(r.Body)
	uniqID := r.Context().Value("uniq_id")
	if uniqID == nil {
		h.log.Error("can not get uniq_id from context")
		h.Error(h.log, w, r, http.StatusInternalServerError, errors.New(""))
		return
	}
	err := h.SessionManager.Delete(uniqID.(string))
	if err != nil {
		h.Error(h.log, w, r, http.StatusInternalServerError, store.DeleteCookieFail)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   uniqID.(string),
		Expires: time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)
	h.Respond(h.log, w, r, http.StatusOK, "successfully logout")
}
