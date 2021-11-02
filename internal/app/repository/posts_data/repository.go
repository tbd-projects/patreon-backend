package repository_posts_data

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_posts_data_repository.go -package=mock_repository -mock_names=Repository=PostsDataRepository . Repository

type Repository interface {
	// Create Errors:
	//		repository_postgresql.UnknownDataFormat
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(postData *models.PostData) (int64, error)

	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(dataID int64) (*models.PostData, error)

	// GetData Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetData(postsId int64) ([]models.PostData, error)

	// Update Errors:
	//		repository_postgresql.UnknownDataFormat
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Update(postData *models.PostData) error

	// ExistsData Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	ExistsData(dataID int64) (bool, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(dataId int64) error
}
