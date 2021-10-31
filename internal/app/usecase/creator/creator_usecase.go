package usecase_creator

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoCreator "patreon/internal/app/repository/creator"
	repoFiles "patreon/internal/app/repository/files"
)

type CreatorUsecase struct {
	repository     repoCreator.Repository
	repositoryFile repoFiles.Repository
}

func NewCreatorUsecase(repository repoCreator.Repository, repositoryFile repoFiles.Repository) *CreatorUsecase {
	return &CreatorUsecase{
		repository:     repository,
		repositoryFile: repositoryFile,
	}
}

// Create Errors:
//		CreatorExist
//		models.IncorrectCreatorNickname
//		models.IncorrectCreatorCategory
//		models.IncorrectCreatorDescription
//		repository_postgresql.IncorrectCategory
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *CreatorUsecase) Create(creator *models.Creator) (int64, error) {
	check, err := usecase.repository.GetCreator(creator.ID)
	if err != nil && err != repository.NotFound {
		return app.InvalidInt, errors.Wrap(err, fmt.Sprintf("METHOD: usecase_creator.Create; "+
			"ERR: error on get creator with ID = %v", creator.ID))
	}
	if check != nil {
		return app.InvalidInt, CreatorExist
	}

	if err = creator.Validate(); err != nil {
		if errors.Is(err, models.IncorrectCreatorCategory) || errors.Is(err, models.IncorrectCreatorNickname) ||
			errors.Is(err, models.IncorrectCreatorDescription) {
			return -1, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Create(creator)
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) GetCreators() ([]models.Creator, error) {
	return usecase.repository.GetCreators()
}

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) GetCreator(id int64) (*models.Creator, error) {
	cr, err := usecase.repository.GetCreator(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("creator with ID = %v not found", id))
	}
	return cr, nil
}

// UpdateCover Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) UpdateCover(data io.Reader, name repoFiles.FileName, id int64) error {
	path, err := usecase.repositoryFile.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return err
	}
	err = usecase.repository.UpdateCover(id, path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" err update cover cretor with id %d", id))
	}
	return nil
}

// UpdateAvatar Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) UpdateAvatar(data io.Reader, name repoFiles.FileName, id int64) error {
	path, err := usecase.repositoryFile.SaveFile(data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	err = usecase.repository.UpdateCover(id, path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" err avatar cover cretor with id %d", id))
	}
	return nil
}
