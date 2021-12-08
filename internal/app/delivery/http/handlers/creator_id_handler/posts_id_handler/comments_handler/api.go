package comments_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_postgresql "patreon/internal/app/repository/comments/postgresql"
)

var codesByErrorsPOST = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	models.InvalidPostId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectPostId, logrus.WarnLevel},
	models.InvalidUserId: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectUserId, logrus.WarnLevel},
	repository_postgresql.CommentAlreadyExist: {
		http.StatusConflict, handler_errors.CommentAlreadyExist, logrus.WarnLevel},
}

var codesByErrorsGET = base_handler.CodeMap{
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
