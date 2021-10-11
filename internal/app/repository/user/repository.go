package repository_user

import (
	"patreon/internal/app/models"
)

type Repository interface {
	// Create Errors:
	// 		LoginAlreadyExist
	// 		NicknameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(*models.User) error

	// FindByLogin Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	FindByLogin(string) (*models.User, error)

	// FindByID Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	FindByID(int64) (*models.User, error)
}
