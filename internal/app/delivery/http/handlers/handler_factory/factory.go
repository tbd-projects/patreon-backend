package handler_factory

import (
	"patreon/internal/app"
	handlers2 "patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/handlers/creator_create_handler"

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
	storage    app.DataStorage
	logger     *logrus.Logger
	urlHandler *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, storage app.DataStorage) *Factory {
	return &Factory{
		storage: storage,
		logger:  logger,
	}
}

func (f *Factory) initAllHandlers() map[int]app.Handler {
	return map[int]app.Handler{
		REGISTER:        handlers2.NewRegisterHandler(f.logger, f.storage),
		LOGIN:           handlers2.NewLoginHandler(f.logger, f.storage),
		LOGOUT:          handlers2.NewLogoutHandler(f.logger, f.storage),
		PROFILE:         handlers2.NewProfileHandler(f.logger, f.storage),
		CREATORS:        handlers2.NewCreatorHandler(f.logger, f.storage),
		CREATOR_WITH_ID: creator_create_handler.NewCreatorCreateHandler(f.logger, f.storage),
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
