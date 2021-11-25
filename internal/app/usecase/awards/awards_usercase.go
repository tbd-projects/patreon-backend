package usecase_awards

import (
	"context"
	"fmt"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	repoAwrds "patreon/internal/app/repository/awards"
	"patreon/internal/microservices/files/delivery/grpc/client"
	repoFiles "patreon/internal/microservices/files/files/repository/files"

	"github.com/pkg/errors"
)

type AwardsUsecase struct {
	repository repoAwrds.Repository
	fileClient client.FileServiceClient
}

func NewAwardsUsecase(repository repoAwrds.Repository, fileClient client.FileServiceClient) *AwardsUsecase {
	return &AwardsUsecase{
		repository: repository,
		fileClient: fileClient,
	}
}

// GetAwards Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AwardsUsecase) GetAwards(creatorId int64) ([]models.Award, error) {
	return usecase.repository.GetAwards(creatorId)
}

// Delete Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (usecase *AwardsUsecase) Delete(id int64) error {
	return usecase.repository.Delete(id)
}

// Update Errors:
// 		repository.NotFound
//		repository_postgresql.NameAlreadyExist
//		repository_postgresql.PriceAlreadyExist
//		models.IncorrectAwardsPrice
//		models.EmptyName
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *AwardsUsecase) Update(awards *models.Award) error {
	if err := awards.Validate(); err != nil {
		if errors.Is(err, models.EmptyName) || errors.Is(err, models.IncorrectAwardsPrice) {
			return err
		}
		return &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Update(awards)
}

// Create Errors:
//		repository_postgresql.NameAlreadyExist
//		repository_postgresql.PriceAlreadyExist
//		models.IncorrectAwardsPrice
//		models.EmptyName
//		app.GeneralError with Errors:
//			app.UnknownError
//			repository.DefaultErrDB
func (usecase *AwardsUsecase) Create(awards *models.Award) (int64, error) {
	if err := awards.Validate(); err != nil {
		if errors.Is(err, models.EmptyName) || errors.Is(err, models.IncorrectAwardsPrice) {
			return app.InvalidInt, err
		}
		return app.InvalidInt, &app.GeneralError{
			Err:         app.UnknownError,
			ExternalErr: errors.Wrap(err, "failed process of validation creator"),
		}
	}

	return usecase.repository.Create(awards)
}

// GetCreatorId Errors:
//  	repository.NotFound
//  	app.GeneralError with Errors
//   		repository.DefaultErrDB
func (usecase *AwardsUsecase) GetCreatorId(awardsId int64) (int64, error) {
	aw, err := usecase.repository.GetByID(awardsId)
	if err != nil {
		return app.InvalidInt, err
	}
	return aw.CreatorId, nil
}

// UpdateCover Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
//			repository_os.ErrorCreate
//   		repository_os.ErrorCopyFile
// 			repository.DefaultErrDB
func (usecase *AwardsUsecase) UpdateCover(data io.Reader, name repoFiles.FileName, awardsId int64) error {
	_, err := usecase.repository.CheckAwards(awardsId)
	if err != nil {
		return err
	}

	path, err := usecase.fileClient.SaveFile(context.Background(), data, name, repoFiles.Image)
	if err != nil {
		return err
	}

	err = usecase.repository.UpdateCover(awardsId, app.LoadFileUrl+path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(" err avatar cover awards with id %d", awardsId))
	}
	return nil
}
