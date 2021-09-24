package handlers

import (
	"encoding/json"
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
	log            *logrus.Logger
	RespondHandler
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		baseHandler: *app.NewHandlerJoiner([]app.Joinable{}, "/profile"),
		log:         logrus.New(),
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
	router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	router.Use(h.authMiddleware.Check)
	h.baseHandler.Join(router)
}
func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(r.Body)
	req := &request{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		h.Error(h.log, w, r, http.StatusUnprocessableEntity, err)
		return
	}
	h.Respond(h.log, w, r, http.StatusOK, "access is allowed")
}
