package posts_id_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
)

var codesByErrorsGET = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.PostNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codesByErrorsDELETE = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
