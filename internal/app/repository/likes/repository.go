package repository_likes

import (
	"patreon/internal/app/models"
)

type Repository interface {
	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(userId int64) (*models.Like, error)

	// Add Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Add(like *models.Like) error

	// Delete Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(likeId int64) error
}
