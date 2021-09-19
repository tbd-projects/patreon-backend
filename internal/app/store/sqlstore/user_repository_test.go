package sqlstore

import (
	"github.com/stretchr/testify/assert"
	"patreon/internal/models"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := TestDB(t, dbUrl)
	defer teardown("users")

	s := New(db)
	u := models.TestUser(t)
	err := s.User().Create(u)
	assert.NoError(t, err)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	db, teardown := TestDB(t, dbUrl)
	defer teardown("users")

	s := New(db)
	login := "mail1999"
	_, err := s.User().FindByLogin(login)
	assert.Error(t, err)

	u := models.TestUser(t)
	u.Login = login

	assert.NoError(t, u.BeforeCreate())
	err = s.User().Create(u)
	assert.NoError(t, err)

	u, err = s.User().FindByLogin(login)
	assert.NotNil(t, u)
	assert.Nil(t, err)
}
