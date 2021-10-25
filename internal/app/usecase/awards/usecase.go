package usecase_awards

import (
	"patreon/internal/app/models"
)

type Usecase interface {
	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Awards, error)

	// GetCreatorId Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreatorId(awardsId int64) (int64, error)

	// Delete Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDBB
	Delete(id int64) error

	// Update Errors:
	// 		repository.NotFound
	//		repository_postgresql.NameAlreadyExist
	//		models.IncorrectAwardsPrice
	//		models.EmptyName
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Update(awards *models.Awards) error

	// Create Errors:
	//		repository_postgresql.NameAlreadyExist
	//		models.IncorrectAwardsPrice
	//		models.EmptyName
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(awards *models.Awards) (int64, error)
}
