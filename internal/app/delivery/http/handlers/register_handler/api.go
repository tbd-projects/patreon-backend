package register_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_user "patreon/internal/app/usecase/user"
	"patreon/internal/models"
)

var codeByError = base_handler.CodeMap{
	repository.NotFound:             {http.StatusNotFound, handler_errors.UserNotFound},
	usecase_user.UserExist:          {http.StatusBadRequest, handler_errors.UserAlreadyExist},
	models.IncorrectEmailOrPassword: {http.StatusBadRequest, handler_errors.IncorrectEmailOrPassword},
	repository.ErrDefaultDB:         {http.StatusInternalServerError, handler_errors.BDError},
}
