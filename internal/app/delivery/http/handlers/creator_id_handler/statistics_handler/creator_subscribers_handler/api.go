package statistics_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase/statistics"
)

var codeByErrorGet = base_handler.CodeMap{
	statistics.CreatorDoesNotExists: {
		http.StatusNotFound, handler_errors.CreatorNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
