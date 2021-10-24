package usecase_awards

import (
	"patreon/internal/app/models"
)

type Usecase interface {
	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Awards, error)

	// Delete Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(id int64) error

	Update(awards *models.Awards, withNameUpdate bool) error

	// Create Errors:
	//		CreatorExist
	//		models.IncorrectCreatorNickname
	//		models.IncorrectCreatorCategory
	//		models.IncorrectCreatorCategoryDescription
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(awards *models.Awards) (int64, error)
}
