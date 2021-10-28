package repository_awards

import "patreon/internal/app/models"

type Repository interface {
	// Create Errors:
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(aw *models.Awards) (int64, error)

	// Update Errors:
	//		repository.NotFound
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Update(aw *models.Awards) error

	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Awards, error)

	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetByID(awardsID int64) (*models.Awards, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(awardsId int64) error
}
