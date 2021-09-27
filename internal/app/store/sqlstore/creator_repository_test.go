package sqlstore

import (
	"patreon/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatorRepository_Create(t *testing.T) {
	db, teardown := TestDB(t, dbUrl)
	defer teardown("users", "creator_profile")

	s := New(db)
	cr := models.TestCreator(t)
	u := models.TestUser(t)
	err := s.User().Create(u)
	assert.NoError(t, err)

	cr.ID = u.ID
	err = s.Creator().Create(cr)
	assert.NoError(t, err)

	scrArray, err := s.Creator().GetCreators()
	assert.NoError(t, err)
	assert.Equal(t, []models.Creator{*cr}, scrArray)

	scr, err := s.Creator().GetCreator(int64(cr.ID))
	assert.NoError(t, err)
	assert.Equal(t, cr, scr)
}
func TestCreatorRepository_GetCreator(t *testing.T) {
	db, teardown := TestDB(t, dbUrl)
	defer teardown("users", "creator_profile")

	s := New(db)
	u := models.TestUser(t)

	err := s.User().Create(u)
	assert.NoError(t, err)

	cr := models.TestCreator(t)
	cr.ID = u.ID
	expected := *cr

	err = s.Creator().Create(cr)
	assert.NoError(t, err)

	get, err := s.Creator().GetCreator(int64(expected.ID))
	assert.NoError(t, err)
	assert.Equal(t, expected, *get)
}
func TestCreatorRepository_GetCreators_AllUsersCreators(t *testing.T) {
	db, teardown := TestDB(t, dbUrl)
	defer teardown("users", "creator_profile")

	s := New(db)
	users := models.TestUsers(t)
	creators := models.TestCreators(t)

	for i, user := range users {
		err := s.User().Create(&user)
		assert.NoError(t, err)
		creators[i].ID = user.ID
		creators[i].Nickname = user.Nickname

	}
	expected := make([]models.Creator, len(creators))
	copy(expected, creators)

	for _, cr := range creators {
		err := s.Creator().Create(&cr)
		assert.NoError(t, err)
	}

	get, err := s.Creator().GetCreators()
	assert.NoError(t, err)
	assert.Equal(t, expected, get)
}
