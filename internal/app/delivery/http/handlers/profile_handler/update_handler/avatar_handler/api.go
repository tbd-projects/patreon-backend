package avatar_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"

	log "github.com/sirupsen/logrus"
)

var codeByError = base_handler.CodeMap{
	app.UnknownError: {http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
}
