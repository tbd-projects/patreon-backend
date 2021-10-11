package usecase_creator

import (
	"fmt"
	"patreon/internal/app"
	repoCreator "patreon/internal/app/repository/creator"
	"patreon/internal/models"

	"github.com/pkg/errors"
)

type CreatorUsecase struct {
	repository repoCreator.Repository
}

func NewCreatorUsecase(repository repoCreator.Repository) *CreatorUsecase {
	return &CreatorUsecase{
		repository: repository,
	}
}
func (usecase *CreatorUsecase) Create(creator *models.Creator) (int64, error) {
	check, err := usecase.repository.GetCreator(creator.ID)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("METHOD: usecase_creator.Create; "+
			"ERR: error on get creator with ID = %v", creator.ID))
	}
	if check != nil {
		return -1, CreatorExist
	}
	if err = creator.Validate(); err != nil {
		return -1, app.GeneralError{
			Err:         err,
			ExternalErr: errors.Wrap(err, "creator data invalid"),
		}
	}
	id, err := usecase.repository.Create(creator)
	if err != nil {
		return -1, err
	}
	return id, nil
}
func (usecase *CreatorUsecase) GetCreators() ([]models.Creator, error) {
	creators, err := usecase.repository.GetCreators()
	if err != nil {
		return nil, err
	}
	return creators, nil
}
func (usecase *CreatorUsecase) GetCreator(id int64) (*models.Creator, error) {
	cr, err := usecase.repository.GetCreator(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("creator with ID = %v not found", id))
	}
	return cr, nil
}
