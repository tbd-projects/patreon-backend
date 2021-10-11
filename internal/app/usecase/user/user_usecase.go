package usercase_user

import (
	"fmt"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoUser "patreon/internal/app/repository/user"
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
			Err: BadEncrypt,
			ExternalErr: err,
		}
	}

	if err = usecase.repository.Create(user); err != nil {
		return -1, err
	}

	return user.ID, nil
}

// Check Errors:
//		IncorrectEmailOrPassword
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *UserUsecase) Check(login string, password string) (int64, error) {
	u, err := usecase.repository.FindByLogin(login)
	if err != nil {
		return -1, err
	}

	if !u.ComparePassword(password) {
		return -1, IncorrectEmailOrPassword
	}
	return u.ID, nil
}
