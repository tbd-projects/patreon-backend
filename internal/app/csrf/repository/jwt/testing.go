package repository_jwt

import (
	"patreon/internal/app/csrf/csrf_models"
	"testing"
	"time"
)

func TestSources(t *testing.T) *csrf_models.TokenSources {
	t.Helper()
	return &csrf_models.TokenSources{
		UserId:      1,
		SessionId:   "session_id",
		ExpiredTime: time.Now().Add(time.Minute),
	}
}
