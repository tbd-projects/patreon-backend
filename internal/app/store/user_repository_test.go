package store

import (
	"github.com/stretchr/testify/assert"
	"patreon/internal/models"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := TestStore(t, dbUrl)
	defer teardown("users")

	u, err := s.User().Create(models.TestUser(t))
	assert.NoError(t, u.BeforeCreate())
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s, teardown := TestStore(t, dbUrl)
	defer teardown("users")

	login := "mail1999"
	_, err := s.User().FindByLogin(login)
	assert.Error(t, err)

	u := models.TestUser(t)
	u.Login = login

	assert.NoError(t, u.BeforeCreate())
	_, err = s.User().Create(u)
	assert.NoError(t, err)

	u, err = s.User().FindByLogin(login)
	assert.NotNil(t, u)
	assert.Nil(t, err)
}
