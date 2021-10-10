package creator_create_handler

import (
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/repository"
	usecase_creator "patreon/internal/app/usecase/creator"
)

var codesByErrors = base_handler.CodeMap{
	usecase_creator.CreatorExist: {http.StatusConflict, handler_errors.ProfileAlreadyExist},
	repository.NotFound:          {http.StatusNotFound, handler_errors.UserNotFound},
	repository.ErrDefaultDB:      {http.StatusInternalServerError, handler_errors.BDError},
}
