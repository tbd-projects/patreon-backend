package creator_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_postgresql "patreon/internal/app/repository/creator/postgresql"
	usecase_creator "patreon/internal/app/usecase/creator"
)

var codesByErrorsGET = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codesByErrorsPOST = base_handler.CodeMap{
	usecase_creator.CreatorExist: {
		http.StatusConflict, handler_errors.CreatorAlreadyExist, logrus.InfoLevel},
	repository.NotFound: {
		http.StatusNotFound, handler_errors.UserNotFound, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	app.UnknownError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	models.IncorrectCreatorCategory: {
		http.StatusUnprocessableEntity, handler_errors.InvalidCategory, logrus.InfoLevel},
	repository_postgresql.IncorrectCategory: {
		http.StatusUnprocessableEntity, handler_errors.InvalidCategory, logrus.InfoLevel},
	models.IncorrectCreatorNickname: {
		http.StatusUnprocessableEntity, handler_errors.InvalidNickname, logrus.InfoLevel},
	models.IncorrectCreatorDescription: {
		http.StatusUnprocessableEntity, handler_errors.InvalidDescription, logrus.InfoLevel},
}