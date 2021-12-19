package payments

import (
	"github.com/sirupsen/logrus"
	"patreon/internal/app/models"
	db_models "patreon/internal/app/models"
	repository_payments "patreon/internal/app/repository/payments"
	push_client "patreon/internal/microservices/push/delivery/client"
)

type PaymentsUsecase struct {
	repository repository_payments.Repository
	pusher     push_client.Pusher
}

func NewPaymentsUsecase(repo repository_payments.Repository, pusher push_client.Pusher) *PaymentsUsecase {
	return &PaymentsUsecase{
		repository: repo,
		pusher:     pusher,
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
func (usecase *PaymentsUsecase) UpdateStatus(log *logrus.Entry, token string, recieveAmount float64) error {
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

	errPush := usecase.pusher.ApplyPayments(token)
	if errPush != nil {
		log.Errorf("Try push new post, and got err %s", errPush)
	}

	return usecase.repository.UpdateStatus(token)
}
