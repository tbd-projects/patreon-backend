package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_create_handler"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/avatar_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler/update_handler/password_handler"
	handlers2 "patreon/internal/app/delivery/http/handlers/register_handler"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ROOT = iota
	REGISTER
	LOGIN
	LOGOUT
	PROFILE
	CREATORS
	CREATOR_WITH_ID
	UPDATE_PASSWORD
	UPDATE_AVATAR
)

type HandlerFactory struct {
	usecaseFactory UsecaseFactory
	logger         *logrus.Logger
	router         *mux.Router
	cors           *app.CorsConfig
	urlHandler     *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, router *mux.Router,
	cors *app.CorsConfig, usecaseFactory UsecaseFactory) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory: usecaseFactory,
		logger:         logger,
		router:         router,
		cors:           cors,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucCreator := f.usecaseFactory.GetCreatorUsecase()
	sManager := f.usecaseFactory.GetSessionManager()
	return map[int]app.Handler{
		REGISTER:        handlers2.NewRegisterHandler(f.logger, f.router, f.cors, sManager, ucUser),
		LOGIN:           login_handler.NewLoginHandler(f.logger, f.router, f.cors, sManager, ucUser),
		LOGOUT:          logout_handler.NewLogoutHandler(f.logger, f.router, f.cors, sManager),
		PROFILE:         profile_handler.NewProfileHandler(f.logger, f.router, f.cors, sManager, ucUser),
		CREATORS:        creator_handler.NewCreatorHandler(f.logger, f.router, f.cors, sManager, ucCreator, ucUser),
		CREATOR_WITH_ID: creator_create_handler.NewCreatorCreateHandler(f.logger, f.router, f.cors, sManager, ucUser, ucCreator),
		UPDATE_PASSWORD: password_handler.NewUpdatePasswordHandler(f.logger, f.router, f.cors, sManager, ucUser),
		UPDATE_AVATAR:   avatar_handler.NewUpdateAvatarHandler(f.logger, f.router, f.cors, sManager, ucUser),
	}
}

func (f *HandlerFactory) GetHandleUrls() *map[string]app.Handler {
	if f.urlHandler != nil {
		return f.urlHandler
	}

	hs := f.initAllHandlers()
	f.urlHandler = &map[string]app.Handler{
		//"/":                     "I am a joke?",
		"/login":                hs[LOGIN],
		"/logout":               hs[LOGOUT],
		"/register":             hs[REGISTER],
		"/user":                 hs[PROFILE],
		"/creators":             hs[CREATORS],
		"/creators/{id:[0-9]+}": hs[CREATOR_WITH_ID],
		"/user/update/password": hs[UPDATE_PASSWORD],
		"/user/update/avatar":   hs[UPDATE_AVATAR],
	}
	return f.urlHandler
}
