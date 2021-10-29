package repository_awards

import "patreon/internal/app/models"

type Repository interface {
	// Create Errors:
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(aw *models.Award) (int64, error)

	// Update Errors:
	//		repository.NotFound
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Update(aw *models.Award) error

	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Award, error)

	// GetByID Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	GetByID(awardsID int64) (*models.Award, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(awardsId int64) error
}
