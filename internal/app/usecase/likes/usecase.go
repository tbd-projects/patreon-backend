package usecase_likes

import (
	"patreon/internal/app/models"
)

type Usecase interface {

	// Add Errors:
	//		usecase_likes.IncorrectAddLike
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Add(like *models.Like) error

	// Delete Errors:
	//		usecase_likes.IncorrectDelLike
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(postId int64, userId int64) error
}
