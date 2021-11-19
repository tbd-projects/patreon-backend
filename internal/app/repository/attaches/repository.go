package repository_attaches

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_attaches_repository.go -package=mock_repository -mock_names=Repository=AttachesRepository . Repository

type Repository interface {
	// Create Errors:
	//		repository_postgresql.UnknownDataFormat
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(postData *models.AttachWithoutLevel) (int64, error)

	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(attachId int64) (*models.AttachWithoutLevel, error)

	// GetAttaches Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAttaches(postsId int64) ([]models.AttachWithoutLevel, error)

	// Update Errors:
	//		repository_postgresql.UnknownDataFormat
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Update(postData *models.AttachWithoutLevel) error

	// ExistsAttach Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	ExistsAttach(attachId ...int64) (bool, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(attachId int64) error

	// ApplyChangeAttaches Errors:
	//		UnknownDataFormat
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
   	ApplyChangeAttaches(postId int64, newAttaches []models.Attach, updatedAttaches []models.Attach) ([]int64, error)
}
