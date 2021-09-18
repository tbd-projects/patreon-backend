package store

import (
	"github.com/stretchr/testify/assert"
	"patreon/internal/models"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := TestStore(t, dbUrl)
	defer teardown("users")

	u, err := s.User().Create(&models.User{
		Login:    "golang@python.js",
		Password: "1234",
		Avatar:   "static/img/avatar.png",
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s, teardown := TestStore(t, dbUrl)
	defer teardown("users")
	login := "i unknown user!"
	_, err := s.User().FindByLogin(login)
	assert.Error(t, err)

	_, err = s.User().Create(&models.User{
		Login:    "golang@python.js",
		Password: "1234",
		Avatar:   "static/img/avatar.png",
	})

	assert.NoError(t, err)

	u, err := s.User().FindByLogin("golang@python.js")
	assert.NotNil(t, u)
	assert.Nil(t, err)
}
