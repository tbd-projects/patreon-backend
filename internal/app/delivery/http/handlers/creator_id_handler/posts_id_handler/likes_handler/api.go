package likes_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_likes "patreon/internal/app/usecase/likes"
)

var codesByErrorsDELETE = base_handler.CodeMap{
	usecase_likes.IncorrectDelLike: {
		http.StatusConflict, handler_errors.LikesAlreadyDel, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}

var codesByErrorsPUT = base_handler.CodeMap{
	usecase_likes.IncorrectAddLike: {
		http.StatusConflict, handler_errors.LikesAlreadyExists, logrus.WarnLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
}
