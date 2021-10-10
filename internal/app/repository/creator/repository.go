package repository_creator

import "patreon/internal/models"

type Repository interface {
	Create(*models.Creator) (int64, error)
	GetCreators() ([]models.Creator, error)
	GetCreator(int64) (*models.Creator, error)
}
