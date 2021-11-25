package login_handler

import (
	"context"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	usecase_user "patreon/internal/app/usecase/user"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"
	"patreon/internal/microservices/auth/sessions/sessions_manager"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"github.com/sirupsen/logrus"
)

type LoginHandler struct {
	sessionClient session_client.AuthCheckerClient
	userUsecase   usecase_user.Usecase
	bh.BaseHandler
}

func NewLoginHandler(log *logrus.Logger, sClient session_client.AuthCheckerClient,
	ucUser usecase_user.Usecase) *LoginHandler {
	h := &LoginHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sClient,
		userUsecase:   ucUser,
	}
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionClient, log).CheckNotAuthorized)
	h.AddMethod(http.MethodPost, h.POST)
	return h
}

// POST Login
// @Summary login user
// @Description login user
// @tags user
// @Accept  json
// @Produce json
// @Param user body http_models.RequestLogin true "Request body for login"
// @Success 200 "Successfully login"
// @Failure 404 {object} http_models.ErrResponse "user not found"
// @Failure 422 {object} http_models.ErrResponse "invalid body in request"
// @Failure 500 {object} http_models.ErrResponse "can not create session", "can not do bd operation"
// @Failure 401 {object} http_models.ErrResponse "incorrect email or password"
// @Failure 418 "User are authorized"
// @Router /login [POST]
func (h *LoginHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestLogin{}
	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil || len(req.Password) == 0 || len(req.Login) == 0 {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	h.Log(r).Debugf("Login : %s, password : %s", req.Login, req.Password)

	id, err := h.userUsecase.Check(req.Login, req.Password)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrors)
		return
	}

	res, err := h.sessionClient.Create(context.Background(), id)
	if err != nil || res.UserID != id {
		h.Log(r).Errorf("Error create session %s", err)
		h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateSession)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    res.UniqID,
		Expires:  time.Now().Add(sessions_manager.ExpiredCookiesTime),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
