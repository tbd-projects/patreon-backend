package posts_update_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

var codesByErrorsPUT = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.PostNotFound, logrus.ErrorLevel},
	models.InvalidAwardsId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectAwardsId, logrus.InfoLevel},
	models.InvalidCreatorId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectCreatorId, logrus.WarnLevel},
	models.EmptyTitle: {
		http.StatusUnprocessableEntity, handler_errors.EmptyTitle, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	app.UnknownError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
