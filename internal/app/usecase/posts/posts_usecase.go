package usecase_posts

import (
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoPosts "patreon/internal/app/repository/posts"
	repoPostsData "patreon/internal/app/repository/posts_data"
)

type PostsUsecase struct {
	repository     repoPosts.Repository
	repositoryData repoPostsData.Repository
}

func NewPostsUsecase(repository repoPosts.Repository, repositoryData repoPostsData.Repository) *PostsUsecase {
	return &PostsUsecase{
		repository:     repository,
		repositoryData: repositoryData,
	}
}

// GetPosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) GetPosts(creatorId int64, userId int64, pag *models.Pagination) ([]models.Post, error) {
	return usecase.repository.GetPosts(creatorId, userId, pag)
}

// GetPost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) GetPost(postId int64, userId int64) (*models.PostWithData, error) {
	post, err := usecase.repository.GetPost(postId, userId)
	if err != nil {
		return nil, err
	}
	res := &models.PostWithData{Post: post, Data: []models.PostData{}}
	res.Data, err = usecase.repositoryData.GetData(postId)
	if err != nil {
		return nil, err
	}
	return res, err
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) Delete(postId int64) error {
	return usecase.repository.Delete(postId)
}

// Update Errors:
// 		repository.NotFound
//		models.InvalidAwardsId
//		models.InvalidCreatorId
//		models.EmptyTitle
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *PostsUsecase) Update(post *models.UpdatePost) error {
	if err := post.Validate(); err != nil {
		if errors.Is(err, models.EmptyTitle) || errors.Is(err, models.InvalidCreatorId) ||
			errors.Is(err, models.InvalidAwardsId) {
			return err
		}
		return &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.UpdatePost(post)
}

// Create Errors:
//		models.InvalidAwardsId
//		models.InvalidCreatorId
//		models.EmptyTitle
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *PostsUsecase) Create(post *models.CreatePost) (int64, error) {
	if err := post.Validate(); err != nil {
		if errors.Is(err, models.EmptyTitle) || errors.Is(err, models.InvalidCreatorId) ||
			errors.Is(err, models.InvalidAwardsId) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Create(post)
}

// GetCreatorId Errors:
//  	repository.NotFound
//  	app.GeneralError with Errors:
//   		repository.DefaultErrDB
func (usecase *PostsUsecase) GetCreatorId(postId int64) (int64, error) {
	aw, err := usecase.repository.GetPostCreator(postId)
	if err != nil {
		return app.InvalidInt, err
	}
	return aw, nil
}