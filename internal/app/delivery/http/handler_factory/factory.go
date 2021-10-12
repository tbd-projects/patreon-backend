package handler_factory

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_create_handler"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	handlers2 "patreon/internal/app/delivery/http/handlers/register_handler"
)

const (
	ROOT = iota
	REGISTER
	LOGIN
	LOGOUT
	PROFILE
	CREATORS
	CREATOR_WITH_ID
)

type HandlerFactory struct {
	usecaseFactory UsecaseFactory
	logger     *logrus.Logger
	urlHandler *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, usecaseFactory UsecaseFactory) *HandlerFactory {
	return &HandlerFactory{
		usecaseFactory: usecaseFactory,
		logger: logger,
	}
}

func (f *HandlerFactory) initAllHandlers() map[int]app.Handler {
	ucUser := f.usecaseFactory.GetUserUsecase()
	ucCreator := f.usecaseFactory.GetCreatorUsecase()
	sManager := f.usecaseFactory.GetSessionManager()
	return map[int]app.Handler{
		REGISTER:        handlers2.NewRegisterHandler(f.logger, sManager, ucUser),
		LOGIN:           login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:          logout_handler.NewLogoutHandler(f.logger, sManager),
		PROFILE:         profile_handler.NewProfileHandler(f.logger, sManager, ucUser),
		CREATORS:        creator_handler.NewCreatorHandler(f.logger, sManager, ucCreator),
		CREATOR_WITH_ID: creator_create_handler.NewCreatorCreateHandler(f.logger, sManager, ucUser, ucCreator),
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
	}
	return f.urlHandler
}
