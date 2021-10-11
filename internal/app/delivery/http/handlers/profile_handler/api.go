package profile_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
)

var codeByError = base_handler.CodeMap{
	repository.NotFound: {http.StatusNotFound, handler_errors.UserNotFound},
}
