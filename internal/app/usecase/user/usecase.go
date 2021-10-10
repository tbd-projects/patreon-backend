package usercase_user

import "patreon/internal/models"

type Usecase interface {
	GetProfile(userID int64) (*models.User, error)
	Create(user *models.User) (int64, error)
	Check(login string, password string) (int64, error)
}
