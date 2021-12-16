package payments

import (
	"patreon/internal/app/models"
	db_models "patreon/internal/app/models"
	repository_payments "patreon/internal/app/repository/payments"
)

type PaymentsUsecase struct {
	repository repository_payments.Repository
}

func NewPaymentsUsecase(repo repository_payments.Repository) *PaymentsUsecase {
	return &PaymentsUsecase{
		repository: repo,
	}
}

// GetUserPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (usecase *PaymentsUsecase) GetUserPayments(userID int64, pag *db_models.Pagination) ([]models.UserPayments, error) {
	userPayments, err := usecase.repository.GetUserPayments(userID, pag)
	if err != nil {
		return nil, err
	}

	return userPayments, nil
}

// GetCreatorPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (usecase *PaymentsUsecase) GetCreatorPayments(creatorID int64, pag *db_models.Pagination) ([]models.CreatorPayments, error) {
	creatorPayments, err := usecase.repository.GetCreatorPayments(creatorID, pag)
	if err != nil {
		return nil, err
	}

	return creatorPayments, nil
}

// UpdateStatus Errors:
//		repository_payments.NotEqualPaymentAmount
//		repository_payments.CountPaymentsByTokenError
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (usecase *PaymentsUsecase) UpdateStatus(token string, recieveAmount float64) error {
	err := usecase.repository.CheckCountPaymentsByToken(token)
	if err != nil {
		return err
	}
	res, err := usecase.repository.GetPaymentByToken(token)
	if err != nil {
		return err
	}
	if res.Amount != recieveAmount {
		return repository_payments.NotEqualPaymentAmount
	}
	return usecase.repository.UpdateStatus(token)
}
