package usercase_user

import (
	"fmt"
	"patreon/internal/app"
	repoUser "patreon/internal/app/repository/user"
	"patreon/internal/models"

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

func (usecase *UserUsecase) GetProfile(userID int64) (*models.User, error) {
	u, err := usecase.repository.FindByID(userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("profile with id %v not found", userID))
	}
	return u, nil
}

func (usecase *UserUsecase) Create(user *models.User) (int64, error) {
	checkUser, err := usecase.repository.FindByLogin(user.Login)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("error on create user with login %v", user.Login))
	}
	if checkUser != nil {
		return -1, UserExist
	}
	if err = user.Validate(); err != nil {
		return -1, app.GeneralError{
			Err:         err,
			ExternalErr: errors.Wrap(err, "user data invalid"),
		}
	}

	if err = user.Encrypt(); err != nil {
		returnedErr := app.GeneralError{
			ExternalErr: err,
		}
		if err == models.EmptyPassword {
			returnedErr.Err = EmptyPassword
		} else {
			returnedErr.Err = BadEncrypt
		}

		return -1, returnedErr
	}

	if err = usecase.repository.Create(user); err != nil {
		return -1, err
	}

	return user.ID, nil
}

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
