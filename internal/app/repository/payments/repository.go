package repository_payments

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_payments_repository.go -package=mock_repository -mock_names=Repository=PaymentsRepository . Repository

type Repository interface {
	// GetUserPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetUserPayments(userID int64) ([]models.Payments, error)
	// GetCreatorPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetCreatorPayments(creatorID int64) ([]models.Payments, error)
}
