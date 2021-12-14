package usecase_comments

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoComments "patreon/internal/app/repository/comments"
	push_client "patreon/internal/microservices/push/delivery/client"
)

type CommentsUsecase struct {
	repository repoComments.Repository
	pusher     push_client.Pusher
}

func NewCommentsUsecase(repository repoComments.Repository, pusher push_client.Pusher) *CommentsUsecase {
	return &CommentsUsecase{
		repository: repository,
		pusher:     pusher,
	}
}

// Create Errors:
//		models.InvalidPostId
//		models.InvalidUserId
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) Create(log *logrus.Entry, cm *models.Comment) (int64, error) {
	if err := cm.Validate(); err != nil {
		if errors.Is(err, models.InvalidPostId) || errors.Is(err, models.InvalidUserId) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}
	commentId, err := usecase.repository.Create(cm)
	errPush := usecase.pusher.NewComment(commentId, cm.PostId, cm.AuthorId)
	if errPush != nil {
		log.Errorf("Try push comment; got error: %s", errPush)
	}
	return commentId, err
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (usecase *CommentsUsecase) Get(commentsId int64) (*models.Comment, error) {
	return usecase.repository.Get(commentsId)
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
