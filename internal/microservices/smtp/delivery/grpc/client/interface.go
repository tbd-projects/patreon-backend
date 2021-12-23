package client

import (
	"context"
)

//go:generate mockgen -destination=mocks/smtp_client_mock.go -package=mock_smtp_client . SmtpServiceClient

type SmtpServiceClient interface {
	Send(ctx context.Context, htmlBody string, emailTo []string) error
	Stop(ctx context.Context) error
}
