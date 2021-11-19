package client

import (
	"context"
	"patreon/internal/microservices/auth/sessions/models"
)

//go:generate mockgen -destination=mocks/auth_checker_mock.go -package=mock_auth_checker . AuthCheckerClient

type AuthCheckerClient interface {
	Check(ctx context.Context, sessionID string) (models.Result, error)
	Create(ctx context.Context, userID int64) (models.Result, error)
	Delete(ctx context.Context, sessionID string) error
}
