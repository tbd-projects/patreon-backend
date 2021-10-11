package usecase_creator

import "patreon/internal/models"

type Usecase interface {
	GetCreators() ([]models.Creator, error)
	GetCreator(id int64) (*models.Creator, error)
	Create(creator *models.Creator) (int64, error)
}
