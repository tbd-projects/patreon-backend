package update_handler

import (
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/sessions"
	usecase_user "patreon/internal/app/usecase/user"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type ProfileUpdateHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    usecase_user.Usecase
	bh.BaseHandler
}

func NewProfileUpdateHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	sManager sessions.SessionsManager, ucUser usecase_user.Usecase) *ProfileUpdateHandler {
	h := &ProfileUpdateHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log, router, cors),
	}
	//h.AddMethod(http.MethodPut, h.)
	//h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, log).Check)
	return h
}
