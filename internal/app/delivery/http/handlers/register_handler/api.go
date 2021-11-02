package register_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_user "patreon/internal/app/repository/user/postgresql"
	useUser "patreon/internal/app/usecase/user"

	"github.com/sirupsen/logrus"
)

var codeByError = base_handler.CodeMap{
	models.EmptyPassword: {
		http.StatusUnprocessableEntity, handler_errors.InvalidBody, logrus.InfoLevel},
	repository_user.LoginAlreadyExist: {
		http.StatusConflict, handler_errors.UserAlreadyExist, logrus.InfoLevel},
	useUser.UserExist: {
		http.StatusConflict, handler_errors.UserAlreadyExist, logrus.InfoLevel},
	repository_user.NicknameAlreadyExist: {
		http.StatusUnprocessableEntity, handler_errors.NicknameAlreadyExist, logrus.InfoLevel},
	models.IncorrectEmailOrPassword: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectLoginOrPassword, logrus.InfoLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
