package nickname_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usercase_user "patreon/internal/app/usecase/user"

	log "github.com/sirupsen/logrus"
)

var codeByErrorPUT = base_handler.CodeMap{
	usercase_user.InvalidOldNickname: {http.StatusUnprocessableEntity, handler_errors.InvalidOldNickname, log.WarnLevel},
	repository.NotFound:              {http.StatusNotFound, handler_errors.UserWithNicknameNotFound, log.WarnLevel},
	usercase_user.NicknameExists:     {http.StatusConflict, handler_errors.NicknameAlreadyExist, log.WarnLevel},
	repository.DefaultErrDB:          {http.StatusInternalServerError, handler_errors.BDError, log.ErrorLevel},
	app.UnknownError:                 {http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
}
