package avatar_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
)

var codeByError = base_handler.CodeMap{
	app.UnknownError: {http.StatusInternalServerError, handler_errors.InternalError},
}
