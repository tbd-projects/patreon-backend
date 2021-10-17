package password_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	usercase_user "patreon/internal/app/usecase/user"
)

var codeByError = base_handler.CodeMap{
	repository.NotFound:                {http.StatusNotFound, handler_errors.UserNotFound},
	usercase_user.IncorrectNewPassword: {http.StatusBadRequest, handler_errors.IncorrectNewPassword},
	models.EmptyPassword:               {http.StatusBadRequest, handler_errors.IncorrectNewPassword},
	repository.DefaultErrDB:            {http.StatusInternalServerError, handler_errors.BDError},
	usercase_user.BadEncrypt:           {http.StatusInternalServerError, handler_errors.InternalError},
	app.UnknownError:                   {http.StatusInternalServerError, handler_errors.InternalError},
	usercase_user.OldPasswordEqualNew:  {http.StatusInternalServerError, handler_errors.BDError},
}
