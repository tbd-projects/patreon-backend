package repository_user

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_user_repository.go -package=mock_repository -mock_names=Repository=UserRepository . Repository

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

	// UpdatePassword Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdatePassword(int64, string) error

	// UpdateAvatar Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdateAvatar(id int64, newAvatar string) error
}
