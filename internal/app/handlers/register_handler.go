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
	"patreon/internal/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	baseHandler    app.HandlerJoiner
	authMiddleware middleware.SessionMiddleware
	router         *mux.Router
	Store          store.Store
	SessionManager sessions.SessionsManager
	RespondHandler
}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/register"),
		RespondHandler: RespondHandler{logrus.New()},
	}
}

func (h *RegisterHandler) SetStore(store store.Store) {
	h.Store = store
}
func (h *RegisterHandler) SetLogger(logger *logrus.Logger) {
	h.log = logger
}
func (h *RegisterHandler) SetSessionManager(manager sessions.SessionsManager) {
	h.SessionManager = manager
	h.authMiddleware = *middleware.NewSessionMiddleware(h.SessionManager, h.log)
}
func (h *RegisterHandler) Join(router *mux.Router) {
	//router.HandleFunc(h.baseHandler.GetUrl(), h.authMiddleware.Check(h)).Methods("POST", "GET")
	//router.Use(gorilla_handlers.CORS(
	//	gorilla_handlers.AllowedOrigins([]string{"http://localhost:3001"}),
	//	gorilla_handlers.AllowedHeaders([]string{"*"}),
	//	gorilla_handlers.AllowCredentials(),
	//	gorilla_handlers.AllowedMethods([]string{"POST", "OPTIONS", "GET"}),
	//))
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST", "GET", "OPTIONS")
	//router.Use(h.authMiddleware.CheckNotAuthorized)

	h.baseHandler.Join(router)
}
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		h.log.Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u := &models.User{
		Login:    req.Login,
		Password: req.Password,
	}

	logUser, _ := json.Marshal(u)
	h.log.Debug("get: ", string(logUser))

	checkUser, _ := h.Store.User().FindByLogin(u.Login)
	if checkUser != nil {
		h.log.Warn(handler_errors.UserAlreadyExist)
		h.Error(w, r, http.StatusConflict, handler_errors.UserAlreadyExist)
		return
	}

	if err := u.Validate(); err != nil {
		h.log.Warnf("Not valid login or password %s", err)
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidBody)
		return
	}

	if err := u.BeforeCreate(); err != nil {
		h.log.Errorf("Error prepare user info %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorPrepareUser)
		return
	}

	if err := h.Store.User().Create(u); err != nil {
		h.log.Errorf("Error create user in bd %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateUser)
		return
	}

	u.MakePrivateDate()
	h.Respond(w, r, http.StatusOK, u)
}
