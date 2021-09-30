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
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST", "GET", "OPTIONS")
	h.baseHandler.Join(router)
}

// Registration
// @Summary create new user
// @Description create new account and get cookies
// @Accept  json
// @Produce json
// @Param user body models.RequestRegistration true "Request body for user registration"
// @Success 201 {object} models.UserResponse "Create user successfully"
// @Failure 400 {object} models.BaseResponse "Invalid body"
// @Failure 409 {object} models.BaseResponse "User already exist"
// @Failure 422 {object} models.BaseResponse "Not valid body"
// @Failure 500 {object} models.BaseResponse "Creation error in base data"
// @Failure 418 "User are authorized"
// @Router /register [POST]
func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.log.Error(err)
		}
	}(r.Body)
	req := &models.RequestRegistration{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil || len(req.Password) == 0 ||
		len(req.Nickname) == 0 || len(req.Login) == 0 {
		h.log.Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u := &models.User{
		Login:    req.Login,
		Password: req.Password,
		Nickname: req.Nickname,
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
