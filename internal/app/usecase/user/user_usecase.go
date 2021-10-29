package usercase_user

import (
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoUser "patreon/internal/app/repository/user"

	"github.com/pkg/errors"
)

type UserUsecase struct {
	repository repoUser.Repository
}

func NewUserUsecase(repository repoUser.Repository) *UserUsecase {
	return &UserUsecase{
		repository: repository,
	}
}

// GetProfile Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *UserUsecase) GetProfile(userID int64) (*models.User, error) {
	u, err := usecase.repository.FindByID(userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("profile with id %v not found", userID))
	}
	return u, nil
}

// Create Errors:
//		models.EmptyPassword
// 		models.IncorrectEmailOrPassword
//		repository_user.LoginAlreadyExist
//		repository_user.NicknameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Create(user *models.User) (int64, error) {
	checkUser, err := usecase.repository.FindByLogin(user.Login)
	if err != nil && err != repository.NotFound {
		return -1, errors.Wrap(err, fmt.Sprintf("error on create user with login %v", user.Login))
	}

	if checkUser != nil {
		return -1, UserExist
	}

	if err = user.Validate(); err != nil {
		if errors.Is(err, models.IncorrectEmailOrPassword) {
			return -1, err
		}
		return -1, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation user"),
		}
	}

	if err = user.Encrypt(); err != nil {
		if errors.Is(err, models.EmptyPassword) {
			return -1, err
		}

		return -1, app.GeneralError{
			Err:         BadEncrypt,
			ExternalErr: err,
		}
	}

	if err = usecase.repository.Create(user); err != nil {
		return -1, err
	}

	return user.ID, nil
}

// Check Errors:
//		models.IncorrectEmailOrPassword
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Check(login string, password string) (int64, error) {
	u, err := usecase.repository.FindByLogin(login)
	if err != nil {
		return -1, err
	}

	if !u.ComparePassword(password) {
		return -1, models.IncorrectEmailOrPassword
	}
	return u.ID, nil
}

// UpdatePassword Errors:
// 		repository.NotFound
//		OldPasswordEqualNew
//		IncorrectNewPassword
//		models.EmptyPassword
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
//			BadEncrypt
//			app.UnknownError
func (usecase *UserUsecase) UpdatePassword(userId int64, newPassword string) error {
	u, err := usecase.GetProfile(userId)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("profile with id %v not found", userId))
	}
	if u.ComparePassword(newPassword) {
		return OldPasswordEqualNew
	}
	u.MakeEmptyPassword()

	u.Password = newPassword
	if err = u.Encrypt(); err != nil {
		if errors.Is(err, models.EmptyPassword) {
			return err
		}
		return app.GeneralError{
			Err:         BadEncrypt,
			ExternalErr: err,
		}
	}
	if err = u.Validate(); err != nil {
		if errors.Is(err, models.IncorrectEmailOrPassword) {
			return IncorrectNewPassword
		}
		return app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation user"),
		}
	}
	err = usecase.repository.UpdatePassword(userId, u.EncryptedPassword)
	return err
}

// UpdateAvatar Errors:
// 		app.GeneralError with Errors
//			app.UnknownError
func (usecase *UserUsecase) UpdateAvatar(userId int64, newAvatar string) error {
	if err := usecase.repository.UpdatePassword(userId, newAvatar); err != nil {
		return app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of update avatar"),
		}
	}
	return nil
}
