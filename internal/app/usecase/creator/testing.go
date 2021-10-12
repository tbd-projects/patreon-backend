package usecase_creator

import (
	"patreon/internal/app/models"
	"testing"
)

func TestCreator(t *testing.T) *models.Creator {
	t.Helper()
	return &models.Creator{
		ID:          1,
		Category:    "podcasts",
		Nickname:    "podcaster2005",
		Description: "blog about IT",
	}
}
