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

	// UpdatePassword Errors:
	// 		repository.NotFound
	//		OldPasswordEqualNew
	//		IncorrectNewPassword
	//		models.EmptyPassword
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	//			BadEncrypt
	//			app.UnknownError
	UpdatePassword(userId int64, newPassword string) error
	// UpdateAvatar Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	//			app.UnknownError
	UpdateAvatar(userId int64, newAvatar string) error
}
