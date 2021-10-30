package posts_data

import (
	"github.com/pkg/errors"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoFiles "patreon/internal/app/repository/files"
	repoPostsData "patreon/internal/app/repository/posts_data"
)

type PostsDataUsecase struct {
	repository      repoPostsData.Repository
	filesRepository repoFiles.Repository
}

func NewPostsDataUsecase(repository repoPostsData.Repository, filesRepository repoFiles.Repository) *PostsDataUsecase {
	return &PostsDataUsecase{
		repository:      repository,
		filesRepository: filesRepository,
	}
}

// GetData Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsDataUsecase) GetData(dataId int64) (*models.PostData, error) {
	return usecase.repository.Get(dataId)
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *PostsDataUsecase) Delete(postId int64) error {
	return usecase.repository.Delete(postId)
}

// LoadImage Errors:
//		models.InvalidPostId
//		models.InvalidType
//		repository_postgresql.UnknownDataFormat
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
func (usecase *PostsDataUsecase) LoadImage(data io.Reader, name repoFiles.FileName, postId int64) (int64, error) {
	path, err := usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return app.InvalidInt, err
	}

	post := &models.PostData{Type: models.Image, Data: app.LoadFileUrl + path, ID: postId}
	if err = post.Validate(); err != nil {
		if errors.Is(err, models.InvalidType) || errors.Is(err, models.InvalidPostId) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}
	return usecase.repository.Create(post)
}

// LoadText Errors:
//		models.InvalidPostId
//		models.InvalidType
//		repository_postgresql.UnknownDataFormat
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *PostsDataUsecase) LoadText(postData *models.PostData) (int64, error) {
	postData.Type = models.Text
	if err := postData.Validate(); err != nil {
		if errors.Is(err, models.InvalidType) || errors.Is(err, models.InvalidPostId) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Create(postData)
}

// UpdateImage Errors:
//		models.InvalidPostId
//		models.InvalidType
//		repository_postgresql.UnknownDataFormat
//		repository.NotFound
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
func (usecase *PostsDataUsecase) UpdateImage(data io.Reader, name repoFiles.FileName, postDataId int64) error {
	if _, err := usecase.repository.Get(postDataId); err != nil {
		return err
	}

	path, err := usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	post := &models.PostData{ID: postDataId, Type: models.Image, Data: app.LoadFileUrl + path}
	if err = post.Validate(); err != nil {
		if errors.Is(err, models.InvalidType) || errors.Is(err, models.InvalidPostId) {
			return err
		}
		return &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}
	return usecase.repository.Update(post)
}

// UpdateText Errors:
//		models.InvalidPostId
//		models.InvalidType
//		repository.NotFound
//		repository_postgresql.UnknownDataFormat
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *PostsDataUsecase) UpdateText(postData *models.PostData) error {
	postData.Type = models.Text
	if err := postData.Validate(); err != nil {
		if errors.Is(err, models.InvalidType) || errors.Is(err, models.InvalidPostId) {
			return err
		}
		return &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Update(postData)
}
