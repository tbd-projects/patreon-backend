package register_handler

import (
	"encoding/json"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	models_respond "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewRegisterHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient,
	ucUser usecase_user.Usecase) *RegisterHandler {
	h := &RegisterHandler{
		sessionClient: sManager,
		userUsecase:   ucUser,
		BaseHandler:   *bh.NewBaseHandler(log),
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
// @Failure 409 {object} models.ErrResponse "user already exist"
// @Failure 422 {object} models.ErrResponse "invalid body in request", "nickname already exist", "incorrect email or password", "incorrect nickname"
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
