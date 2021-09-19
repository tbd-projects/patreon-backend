package teststore

import (
	"github.com/stretchr/testify/assert"
	"patreon/internal/app/store"
	"patreon/internal/models"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st := New()

	u := models.TestUser(t)
	err := st.User().Create(u)
	assert.NoError(t, err)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s := New()
	login := "mail1999"
	_, err := s.User().FindByLogin(login)
	assert.EqualError(t, store.NotFound, err.Error())

	u := models.TestUser(t)
	u.Login = login

	assert.NoError(t, u.BeforeCreate())
	err = s.User().Create(u)
	assert.NoError(t, err)

	u, err = s.User().FindByLogin(login)
	assert.NotNil(t, u)
	assert.Nil(t, err)
}
