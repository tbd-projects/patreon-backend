package posts

import (
	"context"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoAttaches "patreon/internal/app/repository/attaches"
	repoPosts "patreon/internal/app/repository/posts"
	"patreon/internal/microservices/files/delivery/grpc/client"
	repoFiles "patreon/internal/microservices/files/files/repository/files"
	"patreon/pkg/utils"

	"github.com/pkg/errors"
)

type PostsUsecase struct {
	repository      repoPosts.Repository
	repositoryData  repoAttaches.Repository
	filesRepository client.FileServiceClient
	imageConvector  utils.ImageConverter
}

func NewPostsUsecase(repository repoPosts.Repository, repositoryData repoAttaches.Repository,
	fileClient client.FileServiceClient, convector ...utils.ImageConverter) *PostsUsecase {
	conv := utils.ImageConverter(&utils.ConverterToWebp{})
	if len(convector) != 0 {
		conv = convector[0]
	}
	return &PostsUsecase{
		repository:      repository,
		repositoryData:  repositoryData,
		imageConvector:  conv,
		filesRepository: fileClient,
	}
}

// GetPosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) GetPosts(creatorId int64, userId int64,
	pag *models.Pagination, withDraft bool) ([]models.Post, error) {
	return usecase.repository.GetPosts(creatorId, userId, pag, withDraft)
}

// GetAvailablePosts Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) GetAvailablePosts(userID int64, pag *models.Pagination) ([]models.AvailablePost, error) {
	return usecase.repository.GetAvailablePosts(userID, pag)
}

// GetPost Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsUsecase) GetPost(postId int64, userId int64, addView bool) (*models.PostWithAttach, error) {
	post, err := usecase.repository.GetPost(postId, userId, addView)
	if err != nil {
		return nil, err
	}
	res := &models.PostWithAttach{Post: post, Data: []models.AttachWithoutLevel{}}
	res.Data, err = usecase.repositoryData.GetAttaches(postId)
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
//		models.EmptyTitle
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *PostsUsecase) Update(post *models.UpdatePost) error {
	if err := post.Validate(); err != nil {
		if errors.Is(err, models.EmptyTitle) || errors.Is(err, models.InvalidAwardsId) {
			if post.IsDraft && errors.Is(err, models.EmptyTitle) {
				return usecase.repository.UpdatePost(post)
			}
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
			if errors.Is(err, models.EmptyTitle) && post.IsDraft {
				return usecase.repository.Create(post)
			}
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

// LoadCover Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
//			utils.ConvertErr
//  		utils.UnknownExtOfFileName
func (usecase *PostsUsecase) LoadCover(data io.Reader, name repoFiles.FileName, postId int64) error {
	if _, err := usecase.repository.GetPostCreator(postId); err != nil {
		return err
	}

	var err error
	data, name, err = usecase.imageConvector.Convert(context.Background(), data, name)
	if err != nil {
		return errors.Wrap(err, "failed convert to webp of update post cover")
	}

	path, err := usecase.filesRepository.SaveFile(context.Background(), data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	return usecase.repository.UpdateCoverPost(postId, app.LoadFileUrl+path)
}
