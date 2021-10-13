package register_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_user "patreon/internal/app/repository/user"
)

var codeByError = base_handler.CodeMap{
	models.EmptyPassword:                 {http.StatusUnprocessableEntity, handler_errors.InvalidBody},
	repository_user.LoginAlreadyExist:    {http.StatusUnprocessableEntity, handler_errors.UserAlreadyExist},
	repository_user.NicknameAlreadyExist: {http.StatusUnprocessableEntity, handler_errors.NicknameAlreadyExist},
	models.IncorrectEmailOrPassword:      {http.StatusUnprocessableEntity, handler_errors.IncorrectEmailOrPassword},
	repository.DefaultErrDB:              {http.StatusInternalServerError, handler_errors.BDError},
}
