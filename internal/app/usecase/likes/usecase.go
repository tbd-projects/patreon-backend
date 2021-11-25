package usecase_likes

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_likes_usecase.go -package=mock_usecase -mock_names=Usecase=LikesUsecase . Usecase

type Usecase interface {

	// Add Errors:
	//		usecase_likes.IncorrectAddLike
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Add(like *models.Like) (int64, error)

	// Delete Errors:
	//		usecase_likes.IncorrectDelLike
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(postId int64, userId int64) (int64, error)
}
