package repository_payments

import "github.com/pkg/errors"

var (
	CountPaymentsByTokenError = errors.New("payment by token must be once")
	NotEqualPaymentAmount     = errors.New("payment amount from request not equal amount from database")
)
