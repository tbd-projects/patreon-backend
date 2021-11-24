package usecase_awards

import (
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
)

//go:generate mockgen -destination=mocks/mock_awards_usecase.go -package=mock_usecase -mock_names=Usecase=AwardsUsecase . Usecase

type Usecase interface {
	// GetAwards Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAwards(creatorId int64) ([]models.Award, error)

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
	Update(awards *models.Award) error

	// Create Errors:
	//		repository_postgresql.NameAlreadyExist
	//		models.IncorrectAwardsPrice
	//		models.EmptyName
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(awards *models.Award) (int64, error)

	// UpdateCover Errors:
	// 		repository.NotFound
	// 		app.GeneralError with Errors:
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	// 			repository.DefaultErrDB
	UpdateCover(data io.Reader, name repoFiles.FileName, awardsId int64) error
}
