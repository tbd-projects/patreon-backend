package payments

import "patreon/internal/app/models"

//go:generate mockgen -destination=mocks/mock_payments_usecase.go -package=mock_usecase -mock_names=Usecase=PaymentsUsecase . Usecase

type Usecase interface {
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
