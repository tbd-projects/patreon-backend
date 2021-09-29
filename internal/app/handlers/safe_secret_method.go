package handlers

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	gh "patreon/internal/app/handlers/general_handlers"
)

type withHideMethod struct {
	baseHandler *gh.BaseHandler
}

func (h *withHideMethod) Log() *logrus.Logger {
	return h.baseHandler.Log()
}

func (h *withHideMethod) SetLogger(logger *logrus.Logger) {
	h.baseHandler.SetLogger(logger)
}

func (h *withHideMethod) JoinHandlers(joinedHandlers []app.Joinable) {
	h.baseHandler.JoinHandlers(joinedHandlers)
}




