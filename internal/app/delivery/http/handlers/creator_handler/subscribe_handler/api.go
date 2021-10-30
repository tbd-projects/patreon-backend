package subscribe_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_subscribers "patreon/internal/app/usecase/subscribers"

	"github.com/sirupsen/logrus"
)

var codesByErrorsPOST = base_handler.CodeMap{
	usecase_subscribers.SubscriptionAlreadyExists: {
		http.StatusConflict, handler_errors.UserAlreadySubscribe, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
var codesByErrorsDELETE = base_handler.CodeMap{
	usecase_subscribers.SubscriptionsNotFound: {
		http.StatusConflict, handler_errors.SubscribesNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
var codesByErrorsGET = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
