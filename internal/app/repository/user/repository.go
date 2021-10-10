package repository_user

import "patreon/internal/models"

type Repository interface {
	Create(*models.User) error
	FindByLogin(string) (*models.User, error)
	FindByID(int64) (*models.User, error)
}
