package usercase_user

import (
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/app/repository/files"
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
	//		repository_postgresql.LoginAlreadyExist
	//		repository_postgresql.NicknameAlreadyExist
	//		UserExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(user *models.User) (int64, error)

	// Check Errors:
	//		models.IncorrectEmailOrPassword
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Check(login string, password string) (int64, error)

	// UpdatePassword Errors:
	// 		repository.NotFound
	//		OldPasswordEqualNew
	//		IncorrectEmailOrPassword
	//		IncorrectNewPassword
	//		models.EmptyPassword
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	//			BadEncrypt
	//			app.UnknownError
	UpdatePassword(userId int64, oldPassword, newPassword string) error

	// UpdateAvatar Errors:
	// 		app.GeneralError with Errors
	//			app.UnknownError
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	UpdateAvatar(data io.Reader, name repoFiles.FileName, userId int64) error
}
