package usecase_comments

import (
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoComments "patreon/internal/app/repository/comments"
)

type CommentsUsecase struct {
	repository repoComments.Repository
}

func NewCommentsUsecase(repository repoComments.Repository) *CommentsUsecase {
	return &CommentsUsecase{
		repository: repository,
	}
}

// Create Errors:
//		repository_postgresql.CommentAlreadyExist
//		models.InvalidPostId
//		models.InvalidUserId
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) Create(cm *models.Comment) (int64, error) {
	if err := cm.Validate(); err != nil {
		if errors.Is(err, models.InvalidPostId) || errors.Is(err, models.InvalidUserId) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}
	return usecase.repository.Create(cm)
}

// Update Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) Update(cm *models.Comment) error {
	return usecase.repository.Update(cm)
}

// CheckExists Errors:
//		repository_postgresql.CommentAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) CheckExists(commentId int64) error {
	return usecase.repository.CheckExists(commentId)
}

// GetUserComments Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) GetUserComments(userId int64, pag *models.Pagination) ([]models.UserComment, error) {
	return usecase.repository.GetUserComments(userId, pag)
}

// GetPostComments Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) GetPostComments(postId int64, pag *models.Pagination) ([]models.PostComment, error) {
	return usecase.repository.GetPostComments(postId, pag)
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) Delete(commentId int64) error {
	return usecase.repository.Delete(commentId)
}
