package awards

import "patreon/internal/app/models"

type Repository interface {
	// Create Errors:
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(aw *models.Awards) (int64, error)

	// UpdateName Errors:
	//		repository.NotFound
	//		repository_postgresql.NameAlreadyExist
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdateName(awardsId int64, name string) error

	// UpdatePriceDescription Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	UpdatePriceDescription(awardsId int64, price int64, description string) error

	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Awards, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(awardsId int64) (*models.Creator, error)
}
