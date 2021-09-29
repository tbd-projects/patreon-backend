package sqlstore

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/store"
	_ "patreon/internal/app/store"
	"patreon/internal/models"
	_ "patreon/internal/models"
	"regexp"
	"strconv"
	"testing"
	_ "testing"
)

type SuiteUserRepository struct {
	Suite
	store *Store
}

func (s *SuiteUserRepository) SetupSuite() {
	s.InitBD()
	s.store = New(s.DB)
}

func (s *SuiteUserRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}


func (s *SuiteUserRepository) TestUserRepository_Create() {
	u := models.TestUser(s.T())

	u.ID = 1
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar"+
		") VALUES ($1, $2, $3, $4)"+"RETURNING user_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(u.ID)))

	err := s.store.User().Create(u)
	assert.NoError(s.T(), err)
}

func (s *SuiteUserRepository) TestUserRepository_FindByLogin() {

	login := "mail1999"
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, login, encrypted_password from users where login=$1")).
		WithArgs(login).
		WillReturnError(store.NotFound)
	_, err := s.store.User().FindByLogin(login)
	assert.EqualError(s.T(), store.NotFound, err.Error())

	u := models.TestUser(s.T())
	u.Login = login

	assert.NoError(s.T(), u.BeforeCreate())

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, login, encrypted_password from users where login=$1")).
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "encrypted_password"}).
		AddRow(strconv.Itoa(u.ID), u.Login, u.EncryptedPassword))
	var gotten *models.User
	gotten, err = s.store.User().FindByLogin(login)

	assert.NotNil(s.T(), gotten)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gotten.ID, u.ID)
	assert.Equal(s.T(), gotten.Login, u.Login)
	assert.Equal(s.T(), gotten.EncryptedPassword, u.EncryptedPassword)
}

func (s *SuiteUserRepository) TestUserRepository_FindByID() {

	ID := int64(1)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, nickname, avatar from users where user_id=$1")).
		WithArgs(ID).
		WillReturnError(store.NotFound)
	_, err := s.store.User().FindByID(ID)
	assert.EqualError(s.T(), store.NotFound, err.Error())

	u := models.TestUser(s.T())
	u.ID = int(ID)

	assert.NoError(s.T(), u.BeforeCreate())

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, nickname, avatar from users where user_id=$1")).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(strconv.Itoa(u.ID), u.Nickname, u.Avatar))
	var gotten *models.User
	gotten, err = s.store.User().FindByID(ID)

	assert.NotNil(s.T(), gotten)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gotten.ID, u.ID)
	assert.Equal(s.T(), gotten.Nickname, u.Nickname)
	assert.Equal(s.T(), gotten.Avatar, u.Avatar)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}