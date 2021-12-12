package repository_payments

import "github.com/pkg/errors"

var (
	CountPaymentsByTokenError = errors.New("payment by token must be once")
)
