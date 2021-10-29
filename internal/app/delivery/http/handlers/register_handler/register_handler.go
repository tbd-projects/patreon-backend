package register_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"patreon/internal/app"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
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
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).CheckNotAuthorized,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfToken)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)
	req := &models_respond.RequestRegistration{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
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

	u.MakeEmptyPassword()
	h.Respond(w, r, http.StatusCreated, models_respond.IdResponse{ID: id})
}
