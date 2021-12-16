package usecase_info

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_info_usecase.go -package=mock_usecase -mock_names=Usecase=InfoUsecase . Usecase

type Usecase interface {

	// Get Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get() (*models.Info, error)
}
