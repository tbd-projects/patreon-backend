package usecase_creator

import (
	"patreon/internal/app/models"
)

type Usecase interface {
	// GetCreators Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreators() ([]models.Creator, error)

	// GetCreator Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreator(id int64) (*models.Creator, error)

	// Create Errors:
	//		CreatorExist
	//		models.IncorrectCreatorNickname
	//		models.IncorrectCreatorCategory
	//		models.IncorrectCreatorDescription
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(creator *models.Creator) (int64, error)
}
