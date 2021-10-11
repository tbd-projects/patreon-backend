package usercase_user

import (
	"patreon/internal/app/models"
)

type Usecase interface {
	// GetProfile Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetProfile(userID int64) (*models.User, error)

	// Create Errors:
	//		models.EmptyPassword
	// 		models.IncorrectEmailOrPassword
	//		repository_user.LoginAlreadyExist
	//		repository_user.NicknameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(user *models.User) (int64, error)

	// Check Errors:
	//		IncorrectEmailOrPassword
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Check(login string, password string) (int64, error)
}
