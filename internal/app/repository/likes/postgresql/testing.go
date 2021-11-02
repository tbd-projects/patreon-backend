package repository_postgresql

import (
	"patreon/internal/app/models"
	"testing"
)

func TestLike(t *testing.T) *models.Like {
	t.Helper()
	return &models.Like{
		ID:     1,
		UserId: 1,
		Value:  1,
		PostId: 1,
	}
}
