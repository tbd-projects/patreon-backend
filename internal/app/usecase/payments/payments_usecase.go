package payments

import "patreon/internal/app/repository/payments"

type PaymentsUsecase struct {
	repository payments.Repository
}

func NewPaymentsUsecase(repo payments.Repository) *PaymentsUsecase {
	return &PaymentsUsecase{
		repository: repo,
	}
}
