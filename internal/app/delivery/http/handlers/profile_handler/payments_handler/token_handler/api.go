package payments_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	repository_redis "patreon/internal/app/repository/pay_token/redis"

	"github.com/sirupsen/logrus"
)

var codeByErrorGET = base_handler.CodeMap{
	repository_redis.SetError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
