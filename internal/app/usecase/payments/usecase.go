package payments

import "patreon/internal/app/models"

const (
	BaseLimit = 10
)

//go:generate mockgen -destination=mocks/mock_payments_usecase.go -package=mock_usecase -mock_names=Usecase=PaymentsUsecase . Usecase

type Usecase interface {
	// GetUserPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetUserPayments(userID int64) ([]models.UserPayments, error)
	// GetCreatorPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetCreatorPayments(creatorID int64) ([]models.CreatorPayments, error)
}
