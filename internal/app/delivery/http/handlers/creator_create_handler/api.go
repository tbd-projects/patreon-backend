package creator_create_handler

import (
	"net/http"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	usecase_creator "patreon/internal/app/usecase/creator"
)

var codesByErrorsPOST = base_handler.CodeMap{
	usecase_creator.CreatorExist:               {http.StatusConflict, handler_errors.ProfileAlreadyExist},
	repository.NotFound:                        {http.StatusNotFound, handler_errors.UserNotFound},
	repository.DefaultErrDB:                    {http.StatusInternalServerError, handler_errors.BDError},
	app.UnknownError:                           {http.StatusInternalServerError, handler_errors.InternalError},
	models.IncorrectCreatorCategory:            {http.StatusUnprocessableEntity, handler_errors.InvalidCategory},
	models.IncorrectCreatorNickname:            {http.StatusUnprocessableEntity, handler_errors.InvalidNickname},
	models.IncorrectCreatorCategoryDescription: {http.StatusUnprocessableEntity, handler_errors.InvalidCategoryDescription},
}

var codesByErrorsGET = base_handler.CodeMap{
	repository.NotFound:     {http.StatusNotFound, handler_errors.UserNotFound},
	repository.DefaultErrDB: {http.StatusInternalServerError, handler_errors.BDError},
}
