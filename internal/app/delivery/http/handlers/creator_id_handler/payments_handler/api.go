package creator_payments_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"

	"github.com/sirupsen/logrus"
)

var codeByErrorGET = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNoContent, handler_errors.CreatorPaymentsNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
