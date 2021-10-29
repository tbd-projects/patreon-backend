package usecase_likes

import (
	"github.com/pkg/errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoLikes "patreon/internal/app/repository/likes"
)

type LikesUsecase struct {
	repository repoLikes.Repository
}

func NewLikesUsecase(repository repoLikes.Repository) *LikesUsecase {
	return &LikesUsecase{
		repository: repository,
	}
}

// Add Errors:
//		usecase_likes.IncorrectAddLike
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *LikesUsecase) Add(like *models.Like) error {
	_, err := usecase.repository.GetLikeId(like.UserId, like.PostId)
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			like.Value = 1
			return usecase.repository.Add(like)
		}
		return err
	}
	return IncorrectAddLike
}

// Delete Errors:
//		usecase_likes.IncorrectDelLike
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *LikesUsecase) Delete(postId int64, userId int64) error {
	likeId, err := usecase.repository.GetLikeId(userId, postId)
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			return IncorrectDelLike
		}
		return err
	}
	return usecase.repository.Delete(likeId)
}
