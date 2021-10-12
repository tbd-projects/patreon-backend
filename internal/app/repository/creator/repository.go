package repository_creator

import (
	"patreon/internal/app/models"
)

type Repository interface {
	// Create Errors:
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
}
