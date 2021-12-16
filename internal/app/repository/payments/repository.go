package repository_payments

import (
	"patreon/internal/app/models"
	db_models "patreon/internal/app/models"
)

//go:generate mockgen -destination=mocks/mock_payments_repository.go -package=mock_repository -mock_names=Repository=PaymentsRepository . Repository

type Repository interface {
	// GetUserPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetUserPayments(userID int64, pag *db_models.Pagination) ([]models.UserPayments, error)
	// GetCreatorPayments Errors:
	//		repository.NotFound
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetCreatorPayments(creatorID int64, pag *db_models.Pagination) ([]models.CreatorPayments, error)
	// CheckCountPaymentsByToken Errors:
	//		repository_payments.CountPaymentsByTokenError
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	CheckCountPaymentsByToken(token string) error
	// UpdateStatus Errors:
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	UpdateStatus(token string) error
	// GetPaymentByToken Errors:
	//		app.GeneralError with Errors:
	//			repository.DefaultErrDB
	GetPaymentByToken(token string) (models.Payments, error)
}
