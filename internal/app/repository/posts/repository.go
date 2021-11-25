package repository_posts

import (
	"patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_posts_repository.go -package=mock_repository -mock_names=Repository=PostsRepository . Repository

type Repository interface {
	// Create Errors:
	// 		app.GeneralError with Errors
	// 			repository.DefaultErrDB
	Create(post *models.CreatePost) (int64, error)

	// GetPost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPost(postID int64, userId int64, addView bool) (*models.Post, error)

	// GetPostCreator Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPostCreator(postID int64) (int64, error)

	// GetPosts Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPosts(creatorsId int64, userId int64, pag *models.Pagination, withDraft bool) ([]models.Post, error)

	// UpdatePost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdatePost(post *models.UpdatePost) error

	// UpdateCoverPost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	UpdateCoverPost(postId int64, cover string) error

	// Delete Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Delete(postId int64) error
}
