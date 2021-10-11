package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/creator_create_handler"
	"patreon/internal/app/delivery/http/handlers/creator_handler"
	"patreon/internal/app/delivery/http/handlers/login_handler"
	"patreon/internal/app/delivery/http/handlers/logout_handler"
	"patreon/internal/app/delivery/http/handlers/profile_handler"
	handlers2 "patreon/internal/app/delivery/http/handlers/register_handler"
	repository_creator "patreon/internal/app/repository/creator"
	repository_user "patreon/internal/app/repository/user"
	"patreon/internal/app/sessions/repository"
	"patreon/internal/app/sessions/sessions_manager"
	usecase_creator "patreon/internal/app/usecase/creator"
	usercase_user "patreon/internal/app/usecase/user"

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
)

type Factory struct {
	connections app.ExpectedConnections
	//storage    app.DataStorage
	logger     *logrus.Logger
	urlHandler *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, connections app.ExpectedConnections) *Factory {
	return &Factory{
		connections: connections,
		//storage: storage,
		logger: logger,
	}
}

func (f *Factory) initAllHandlers() map[int]app.Handler {
	ucUser := usercase_user.NewUserUsecase(repository_user.NewUserRepository(f.connections.SqlConnection))
	ucCreator := usecase_creator.NewCreatorUsecase(repository_creator.NewCreatorRepository(f.connections.SqlConnection))
	sManager := sessions_manager.NewSessionManager(repository.NewRedisRepository(f.connections.RedisPool, f.logger))
	return map[int]app.Handler{
		REGISTER:        handlers2.NewRegisterHandler(f.logger, sManager, ucUser),
		LOGIN:           login_handler.NewLoginHandler(f.logger, sManager, ucUser),
		LOGOUT:          logout_handler.NewLogoutHandler(f.logger, sManager),
		PROFILE:         profile_handler.NewProfileHandler(f.logger, sManager, ucUser),
		CREATORS:        creator_handler.NewCreatorHandler(f.logger, sManager, ucCreator),
		CREATOR_WITH_ID: creator_create_handler.NewCreatorCreateHandler(f.logger, sManager, ucUser, ucCreator),
	}
}

func (f *Factory) GetHandleUrls() *map[string]app.Handler {
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
