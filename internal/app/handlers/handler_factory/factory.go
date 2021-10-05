package handler_factory

import (
	"patreon/internal/app"
	"patreon/internal/app/handlers"

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
	storage    *app.DataStorage
	logger     *logrus.Logger
	urlHandler *map[string]app.Handler
}

func NewFactory(logger *logrus.Logger, storage *app.DataStorage) *Factory {
	return &Factory{
		storage: storage,
		logger:  logger,
	}
}

func (f *Factory) initAllHandlers() map[int]app.Handler {
	return map[int]app.Handler{
		REGISTER:        handlers.NewRegisterHandler(f.logger, f.storage),
		LOGIN:           handlers.NewLoginHandler(f.logger, f.storage),
		LOGOUT:          handlers.NewLogoutHandler(f.logger, f.storage),
		PROFILE:         handlers.NewProfileHandler(f.logger, f.storage),
		CREATORS:        handlers.NewCreatorHandler(f.logger, f.storage),
		CREATOR_WITH_ID: handlers.NewCreatorCreateHandler(f.logger, f.storage),
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
		"/profile":              hs[PROFILE],
		"/creators":             hs[CREATORS],
		"/creators/{id:[0-9]+}": hs[CREATOR_WITH_ID],
	}
	return f.urlHandler
}
