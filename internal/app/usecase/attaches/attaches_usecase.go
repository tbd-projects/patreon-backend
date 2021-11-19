package attaches

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoAttaches "patreon/internal/app/repository/attaches"
	repoFiles "patreon/internal/app/repository/files"
)

type AttachesUsecase struct {
	repository      repoAttaches.Repository
	filesRepository repoFiles.Repository
}

func NewAttachesUsecase(repository repoAttaches.Repository, filesRepository repoFiles.Repository) *AttachesUsecase {
	return &AttachesUsecase{
		repository:      repository,
		filesRepository: filesRepository,
	}
}

// GetAttach Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AttachesUsecase) GetAttach(attachId int64) (*models.AttachWithoutLevel, error) {
	return usecase.repository.Get(attachId)
}


func (usecase *AttachesUsecase) processingValidateErrorAttach(err error) error {
	if !(errors.Is(err, models.IncorrectType) || errors.Is(err, models.IncorrectAttachId) ||
		errors.Is(err, models.IncorrectLevel)) {
		err = &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation attach"),
		}
	}
	return err
}

// UpdateAttach Errors:
//		repository.NotFound
//		models.IncorrectType
//  	models.IncorrectAttachId
//      models.IncorrectLevel
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AttachesUsecase) checkAttach(newAttach []models.Attach, updatedAttach []models.Attach) error {
	var err error
	for _, att := range newAttach {
		if err = att.Validate(); err != nil {
			err = usecase.processingValidateErrorAttach(err)
			break
		}
	}

	if err != nil {
		return err
	}

	var checkIds []int64
	for _, att := range updatedAttach {
		if err = att.Validate(); err != nil {
			err = usecase.processingValidateErrorAttach(err)
		}

		if att.Id <= 0 {
			err = models.IncorrectAttachId
			break
		}
		checkIds = append(checkIds, att.Id)
	}

	if err != nil {
		return err
	}

	_, err = usecase.repository.ExistsAttach(checkIds...)
	return err
}

// UpdateAttach Errors:
//		repository.NotFound
//		repository_postgresql.UnknownDataFormat
//		models.IncorrectType
//  	models.IncorrectAttachId
//      models.IncorrectLevel
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AttachesUsecase) UpdateAttach(postId int64,
	newAttaches []models.Attach, updatedAttaches []models.Attach) ([]int64, error) {
	if err := usecase.checkAttach(newAttaches, updatedAttaches); err != nil {
		return nil, err
	}

	res, err := usecase.repository.ApplyChangeAttaches(postId, newAttaches, updatedAttaches)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("err with add attaches %d", postId))
	}

	return res, nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AttachesUsecase) Delete(postId int64) error {
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
func (usecase *AttachesUsecase) LoadImage(data io.Reader, name repoFiles.FileName, postId int64) (int64, error) {
	path, err := usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return app.InvalidInt, err
	}

	post := &models.AttachWithoutLevel{Type: models.Image, Value: app.LoadFileUrl + path, PostId: postId}
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
func (usecase *AttachesUsecase) LoadText(postData *models.AttachWithoutLevel) (int64, error) {
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
func (usecase *AttachesUsecase) UpdateImage(data io.Reader, name repoFiles.FileName, postDataId int64) error {
	if _, err := usecase.repository.ExistsAttach(postDataId); err != nil {
		return err
	}

	path, err := usecase.filesRepository.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	post := &models.AttachWithoutLevel{ID: postDataId, Type: models.Image, Value:  app.LoadFileUrl + path}
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
func (usecase *AttachesUsecase) UpdateText(postData *models.AttachWithoutLevel) error {
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
