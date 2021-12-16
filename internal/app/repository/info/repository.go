package repository_info

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_info_repository.go -package=mock_repository -mock_names=Repository=InfoRepository . Repository

type Repository interface {
	// Get Errors:
	// 		app.GeneralError with Errors:
	// 			repository.DefaultErrDB
	Get() (*models.Info, error)
}
