package repository_creator

import (
	"patreon/internal/app/models"
)

type Repository interface {
	// Create Errors:
	//		repository_postgresql.IncorrectCategory
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(*models.Creator) (int64, error)

	// GetCreators Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreators() ([]models.Creator, error)

	// GetCreator Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreator(int64) (*models.Creator, error)

	// ExistsCreator Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	ExistsCreator(creatorId int64) (bool, error)

	// UpdateAvatar Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateAvatar(creatorId int64, avatar string) error

	// UpdateCover Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateCover(creatorId int64, cover string) error
}
