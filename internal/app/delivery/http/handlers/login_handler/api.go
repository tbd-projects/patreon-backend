package login_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	"patreon/internal/models"
)

var codesByErrors = base_handler.CodeMap{
	repository.NotFound:             {http.StatusNotFound, handler_errors.UserNotFound},
	repository.ErrDefaultDB:         {http.StatusInternalServerError, handler_errors.BDError},
	models.IncorrectEmailOrPassword: {http.StatusUnauthorized, handler_errors.IncorrectEmailOrPassword},
	models.InternalError:            {http.StatusInternalServerError, handler_errors.InternalError},
}
