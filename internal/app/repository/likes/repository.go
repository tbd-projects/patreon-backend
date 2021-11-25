package repository_likes

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_likes_repository.go -package=mock_repository -mock_names=Repository=LikesRepository . Repository

type Repository interface {
	// Get Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get(userId int64) (*models.Like, error)

	// GetLikeId Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetLikeId(userId int64, postId int64) (int64, error)

	// Add Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Add(like *models.Like) (int64, error)

	// Delete Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(likeId int64) (int64, error)
}
