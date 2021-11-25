package usecase_creator

import (
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
)

//go:generate mockgen -destination=mocks/mock_creator_usecase.go -package=mock_usecase -mock_names=Usecase=CreatorUsecase . Usecase

type Usecase interface {
	// GetCreators Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreators() ([]models.Creator, error)

	// SearchCreators Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	SearchCreators(pag *models.Pagination, searchString string, categories ...string) ([]models.Creator, error)

	// GetCreator Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetCreator(id int64, userId int64) (*models.CreatorWithAwards, error)

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
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	// 			repository.DefaultErrDB
	UpdateCover(data io.Reader, name repoFiles.FileName, id int64) error

	// UpdateAvatar Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	// 			repository.DefaultErrDB
	UpdateAvatar(data io.Reader, name repoFiles.FileName, id int64) error
}
