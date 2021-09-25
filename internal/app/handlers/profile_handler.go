package handlers

import (
	"errors"
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProfileHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	router         *mux.Router
	Store          store.Store
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/profile"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *ProfileHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *ProfileHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *ProfileHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *ProfileHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("POST", "GET")

	//router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	//router.Use(h.authMiddleware.Check)
	h.baseHandler.Join(router)
}
func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(r.Body)
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.log.Error("can not get user_id from context")
		h.Error(h.log, w, r, http.StatusInternalServerError, errors.New(""))
		return
	}

	u, err := h.Store.User().FindByID(userID.(int64))
	if err != nil {
		h.log.Errorf("get: %s err:%s can not get user from db", u, err.Error())
		h.Error(h.log, w, r, http.StatusServiceUnavailable, store.GetProfileFail)
		return
	}
	h.log.Infof("get profile %s", u)
	h.Respond(h.log, w, r, http.StatusOK, u)
}
