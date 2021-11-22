package upd_audio_attach_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_postgresql "patreon/internal/app/repository/attaches/postgresql"
	repository_os "patreon/internal/app/repository/files/os"

	"github.com/sirupsen/logrus"
)

var codeByErrorPUT = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.AttachNotFound, logrus.ErrorLevel},
	repository_postgresql.UnknownDataFormat: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectDataType, logrus.WarnLevel},
	models.InvalidType: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectDataType, logrus.WarnLevel},
	models.InvalidPostId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectPostId, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	app.UnknownError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_os.ErrorCopyFile: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
	repository_os.ErrorCreate: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
