package posts

import (
	"github.com/sirupsen/logrus"
	"io"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
)

const (
	BaseLimit = 10
	EmptyUser = -2
)

//go:generate mockgen -destination=mocks/mock_posts_usecase.go -package=mock_usecase -mock_names=Usecase=PostsUsecase . Usecase

type Usecase interface {

	// GetAvailablePosts Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetAvailablePosts(userID int64, pag *models.Pagination) ([]models.AvailablePost, error)
	// GetPosts Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPosts(creatorId int64, userId int64, pag *models.Pagination, withDraft bool) ([]models.Post, error)

	// GetPost Errors:
	//		repository.NotFound
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	GetPost(postId int64, userId int64, addView bool) (*models.PostWithAttach, error)

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
	Update(log *logrus.Entry, post *models.UpdatePost) error

	// Create Errors:
	//		models.InvalidAwardsId
	//		models.InvalidCreatorId
	//		models.EmptyTitle
	//		app.GeneralError with Errors:
	//			app.UnknownError
	//			repository.DefaultErrDB
	Create(log *logrus.Entry, post *models.CreatePost) (int64, error)

	// GetCreatorId Errors:
	//  	repository.NotFound
	//  	app.GeneralError with Errors:
	//   		repository.DefaultErrDB
	GetCreatorId(postId int64) (int64, error)

	// LoadCover Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	//			repository_os.ErrorCreate
	//   		repository_os.ErrorCopyFile
	LoadCover(data io.Reader, name repoFiles.FileName, postId int64) error
}
