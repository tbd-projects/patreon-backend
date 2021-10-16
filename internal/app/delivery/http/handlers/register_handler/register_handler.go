package register_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	models_respond "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	usecase_user "patreon/internal/app/usecase/user"

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
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).CheckNotAuthorized)
	return h
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
func (h *RegisterHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	req := &models_respond.RequestRegistration{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil || len(req.Password) == 0 ||
		len(req.Nickname) == 0 || len(req.Login) == 0 {
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

	u.MakePrivateDate()
	h.Respond(w, r, http.StatusOK, id)
}
