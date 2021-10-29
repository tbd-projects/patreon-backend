package aw_other_update_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_postgresql "patreon/internal/app/repository/awards/postgresql"
)

var codesByErrorsPUT = base_handler.CodeMap{
	repository.NotFound: {
		http.StatusNotFound, handler_errors.AwardNotFound, logrus.ErrorLevel},
	repository.DefaultErrDB: {
		http.StatusInternalServerError, handler_errors.BDError, logrus.ErrorLevel},
	repository_postgresql.NameAlreadyExist: {
		http.StatusConflict, handler_errors.AwardsAlreadyExists, logrus.InfoLevel},
	models.EmptyName: {
		http.StatusUnprocessableEntity, handler_errors.EmptyName, logrus.WarnLevel},
	models.IncorrectAwardsPrice: {
		http.StatusUnprocessableEntity, handler_errors.IncorrectPrice, logrus.WarnLevel},
	app.UnknownError: {
		http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
