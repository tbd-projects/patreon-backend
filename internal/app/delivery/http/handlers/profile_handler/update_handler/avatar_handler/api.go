package avatar_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	repository_os "patreon/internal/app/repository/files/os"

	log "github.com/sirupsen/logrus"
)

var codeByError = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, log.ErrorLevel},
	repository.NotFound: {
		http.StatusUnprocessableEntity, handler_errors.UserNotFound, log.WarnLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, log.ErrorLevel},
}
