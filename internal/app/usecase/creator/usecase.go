package usecase_creator

import (
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/app/repository/files"
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

	// UpdateCover Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateCover(data io.Reader, name repoFiles.FileName, id int64) error

	// UpdateAvatar Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateAvatar(data io.Reader, name repoFiles.FileName, id int64) error
}
