package usecase_creator

import (
	"fmt"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repoCreator "patreon/internal/app/repository/creator"
)

type CreatorUsecase struct {
	repository repoCreator.Repository
}

func NewCreatorUsecase(repository repoCreator.Repository) *CreatorUsecase {
	return &CreatorUsecase{
		repository: repository,
	}
}

// Create Errors:
//		CreatorExist
//		models.IncorrectCreatorNickname
//		models.IncorrectCreatorCategory
//		models.IncorrectCreatorCategoryDescription
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *CreatorUsecase) Create(creator *models.Creator) (int64, error) {
	check, err := usecase.repository.GetCreator(creator.ID)
	if err != nil && err != repository.NotFound {
		return -1, errors.Wrap(err, fmt.Sprintf("METHOD: usecase_creator.Create; "+
			"ERR: error on get creator with ID = %v", creator.ID))
	}
	if check != nil {
		return -1, CreatorExist
	}

	if err = creator.Validate(); err != nil {
		if errors.Is(err, models.IncorrectCreatorCategory) || errors.Is(err, models.IncorrectCreatorNickname) ||
			errors.Is(err, models.IncorrectCreatorCategoryDescription) {
			return -1, err
		}
		return -1, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	id, err := usecase.repository.Create(creator)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *CreatorUsecase) GetCreators() ([]models.Creator, error) {
	creators, err := usecase.repository.GetCreators()
	if err != nil {
		return nil, err
	}
	return creators, nil
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
