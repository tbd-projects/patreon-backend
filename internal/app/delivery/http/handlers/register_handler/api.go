package register_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	useUser "patreon/internal/app/usecase/user"
)

var codeByError = base_handler.CodeMap{
	repository.NotFound:              {http.StatusNotFound, handler_errors.UserNotFound},
	useUser.UserExist:                {http.StatusBadRequest, handler_errors.UserAlreadyExist},
	useUser.IncorrectEmailOrPassword: {http.StatusBadRequest, handler_errors.IncorrectEmailOrPassword},
	repository.DefaultErrDB:          {http.StatusInternalServerError, handler_errors.BDError},
}
