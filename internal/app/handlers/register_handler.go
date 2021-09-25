package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"patreon/internal/app"
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
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST")
	//router.Use(h.authMiddleware.CheckNotAuthorized)
	h.baseHandler.Join(router)
}

// Registration
// @Summary create new user
// @Description create new account and get cookies
// @Accept  json
// @Produce json
// @Success 201 {object} models.UserResponse "Create user successfully"
// @Failure 400 {object} models.BaseResponse "Invalid body"
// @Failure 409 {object} models.BaseResponse "User already exist"
// @Failure 422 {object} models.BaseResponse "Not valid body"
// @Failure 500 {object} models.BaseResponse "Creation error in base data"
// @Failure 418 "User are authorized"
// @Router /register [POST]
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
		h.Error(w, r, http.StatusUnprocessableEntity, errors.New("store.InvalidBody"))
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
		h.Error(w, r, http.StatusConflict, store.UserAlreadyExist)
		return
	}
	//TODO(Дима Варин): Разделить валидацию и запрос в бд Если невалиден то кидать 400, если ошибка бд то кидать 500
	if err := h.Store.User().Create(u); err != nil {
		h.Error(w, r, http.StatusBadRequest, err)
		return
	}

	u.MakePrivateDate()
	h.Respond(w, r, http.StatusOK, u)
}
