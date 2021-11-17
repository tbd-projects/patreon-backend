package client

import (
	"context"
	"patreon/internal/microservices/auth/sessions/models"
)

type AuthCheckerClient interface {
	Check(ctx context.Context, sessionID string) (models.Result, error)
	Create(ctx context.Context, userID int64) (models.Result, error)
	Delete(ctx context.Context, sessionID string) error
}
