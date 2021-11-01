package password_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	usercase_user "patreon/internal/app/usecase/user"

	log "github.com/sirupsen/logrus"
)

var codeByError = base_handler.CodeMap{
	repository.NotFound:                {http.StatusNotFound, handler_errors.UserNotFound, log.WarnLevel},
	usercase_user.IncorrectNewPassword: {http.StatusBadRequest, handler_errors.IncorrectNewPassword, log.InfoLevel},
	models.EmptyPassword:               {http.StatusBadRequest, handler_errors.IncorrectNewPassword, log.InfoLevel},
	repository.DefaultErrDB:            {http.StatusInternalServerError, handler_errors.BDError, log.ErrorLevel},
	usercase_user.BadEncrypt:           {http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	app.UnknownError:                   {http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	usercase_user.OldPasswordEqualNew:  {http.StatusBadRequest, handler_errors.BDError, log.InfoLevel},
}
