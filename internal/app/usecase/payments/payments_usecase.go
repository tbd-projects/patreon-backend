package payments

import (
	"patreon/internal/app/models"
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
func (usecase *PaymentsUsecase) GetUserPayments(userID int64) ([]models.UserPayments, error) {
	userPayments, err := usecase.repository.GetUserPayments(userID)
	if err != nil {
		return nil, err
	}

	return userPayments, nil
}

// GetCreatorPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (usecase *PaymentsUsecase) GetCreatorPayments(creatorID int64) ([]models.CreatorPayments, error) {
	creatorPayments, err := usecase.repository.GetCreatorPayments(creatorID)
	if err != nil {
		return nil, err
	}

	return creatorPayments, nil
}
