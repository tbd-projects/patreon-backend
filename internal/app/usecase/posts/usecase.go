package usecase_posts

import (
	"patreon/internal/app/models"
)

const (
	BaseLimit = 10
	EmptyUser = -2
)

type Usecase interface {

	// GetPosts Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPosts(creatorId int64, userId int64, pag *models.Pagination) ([]models.Post, error)

	// GetPost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPost(postId int64, userId int64) (*models.PostWithData, error)

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(postId int64) error

	// Update Errors:
	// 		repository.NotFound
	//		models.InvalidAwardsId
	//		models.InvalidCreatorId
	//		models.EmptyTitle
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Update(post *models.UpdatePost) error

	// Create Errors:
	//		models.InvalidAwardsId
	//		models.InvalidCreatorId
	//		models.EmptyTitle
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(post *models.CreatePost) (int64, error)

	// GetCreatorId Errors:
	//  	repository.NotFound
	//  	app.GeneralError with Errors:
	//   		repository.DefaultErrDB
	GetCreatorId(postId int64) (int64, error)
}
