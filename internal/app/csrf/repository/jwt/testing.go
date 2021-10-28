package repository_jwt

import (
	"patreon/internal/app/csrf/models"
	"testing"
	"time"
)

func TestSources(t *testing.T) *models.TokenSources {
	t.Helper()
	return &models.TokenSources{
		UserId:      1,
		SessionId:   "session_id",
		ExpiredTime: time.Now().Add(time.Minute),
	}
}
