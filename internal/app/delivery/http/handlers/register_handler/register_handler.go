package register_handler

import (
	"encoding/json"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	models_respond "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	"patreon/internal/app/sessions"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewRegisterHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig, sManager sessions.SessionsManager,
	ucUser usecase_user.Usecase) *RegisterHandler {
	h := &RegisterHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST Registration
// @Summary create new user
// @Description create new account and get cookies
// @Accept  json
// @Produce json
// @Param user body models.RequestRegistration true "Request body for user registration"
// @Success 201 {object} models.IdResponse "Create user successfully"
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 409 {object} models.ErrResponse "user already exist"
// @Failure 422 {object} models.ErrResponse "nickname already exist"
// @Failure 422 {object} models.ErrResponse "incorrect email or password"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 418 "User are authorized"
// @Router /register [POST]
func (h *RegisterHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &models_respond.RequestRegistration{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u := &models.User{
		Login:    req.Login,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	logUser, _ := json.Marshal(u)
	h.Log(r).Debug("get: ", string(logUser))

	id, err := h.userUsecase.Create(u)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	u.MakeEmptyPassword()
	h.Respond(w, r, http.StatusCreated, models_respond.IdResponse{ID: id})
}
