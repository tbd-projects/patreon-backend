package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/handlers/handler_errors"
	"patreon/internal/app/sessions/middleware"
	"patreon/internal/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	baseHandler    app.HandlerJoiner
	dataStorage    *app.DataStorage
	authMiddleware middleware.SessionMiddleware
	RespondHandler
}

func NewLoginHandler(storage *app.DataStorage) *LoginHandler {
	h := &LoginHandler{
		baseHandler:    *app.NewHandlerJoiner([]app.Joinable{}, "/login"),
		dataStorage:    storage,
		RespondHandler: RespondHandler{logrus.New()},
	}
	if storage != nil {
		h.authMiddleware = *middleware.NewSessionMiddleware(h.dataStorage.SessionManager, h.log)
	}

	return h
}

func (h *LoginHandler) Join(router *mux.Router) {
	router.Handle(h.baseHandler.GetUrl(), h.authMiddleware.CheckNotAuthorized(h)).Methods("POST", "OPTIONS")
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
			h.log.Fatal(err)
		}
	}(r.Body)

	req := &models.RequestLogin{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(req); err != nil ||
		len(req.Login) == 0 || len(req.Password) == 0 {
		h.log.Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	u, err := h.dataStorage.Store.User().FindByLogin(req.Login)
	h.log.Debugf("Login : %s, password : %s", req.Login, req.Password)
	if err != nil || !u.ComparePassword(req.Password) {
		h.log.Warnf("Fail get user or compare password %s", err)
		h.Error(w, r, http.StatusUnauthorized, handler_errors.IncorrectEmailOrPassword)
		return
	}

	res, err := h.dataStorage.SessionManager.Create(int64(u.ID))
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
