package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	router         *mux.Router
	Store          store.Store
	SessionManager sessions.SessionsManager
	log            *logrus.Logger
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		baseHandler: *app.NewHandlerJoiner([]app.Joinable{}, "/login"),
		log:         logrus.New(),
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
	router.HandleFunc(h.baseHandler.GetUrl(), h.ServeHTTP).Methods("POST", "GET")
	router.Use(h.authMiddleware.Check)
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
			log.Fatal(err)
		}
	}(r.Body)

	req := &request{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil {
		h.Error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	u, err := h.Store.User().FindByLogin(req.Login)

	if err != nil || !u.ComparePassword(req.Password) {
		h.Error(w, r, http.StatusUnauthorized, store.IncorrectEmailOrPassword)
		return
	}
	h.Respond(w, r, http.StatusOK, "successfully login")
}
func (h *LoginHandler) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.Respond(w, r, code, map[string]string{"error": err.Error()})
}
func (h *LoginHandler) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	encoder := json.NewEncoder(w)
	w.WriteHeader(code)
	if data != nil {
		err := encoder.Encode(data)
		if err != nil {
			h.log.Error(err)
		}
	}
	logUser, _ := json.Marshal(data)
	h.log.Info("Respond data: ", string(logUser))
}
