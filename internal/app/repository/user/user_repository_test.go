package repository_user

import (
	"patreon/internal/app/repository"
	"patreon/internal/models"
	_ "patreon/internal/models"
	"regexp"
	"strconv"
	"testing"
	_ "testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SuiteUserRepository struct {
	repository.Suite
	repo *UserRepository
}

func (s *SuiteUserRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewUserRepository(s.DB)
}

func (s *SuiteUserRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteUserRepository) TestUserRepository_Create() {
	u := models.TestUser(s.T())

	u.ID = 1
	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar"+
		") VALUES ($1, $2, $3, $4)"+"RETURNING user_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(u.ID)))

	err := s.repo.Create(u)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar"+
		") VALUES ($1, $2, $3, $4)"+"RETURNING user_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).WillReturnError(repository.BDError)

	err = s.repo.Create(u)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.BDError, err)
}

func (s *SuiteUserRepository) TestUserRepository_FindByLogin() {

	login := "mail1999"
	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, login, encrypted_password from users where login=$1")).
		WithArgs(login).
		WillReturnError(repository.BDError)
	_, err := s.repo.FindByLogin(login)
	assert.EqualError(s.T(), NotFound, err.Error())

	u := models.TestUser(s.T())
	u.Login = login

	assert.NoError(s.T(), u.Encrypt())

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, login, encrypted_password from users where login=$1")).
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "encrypted_password"}).
			AddRow(strconv.Itoa(u.ID), u.Login, u.EncryptedPassword))
	var gotten *models.User
	gotten, err = s.repo.FindByLogin(login)

	assert.NotNil(s.T(), gotten)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gotten.ID, u.ID)
	assert.Equal(s.T(), gotten.Login, u.Login)
	assert.Equal(s.T(), gotten.EncryptedPassword, u.EncryptedPassword)
}

func (s *SuiteUserRepository) TestUserRepository_FindByID() {

	ID := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, nickname, avatar from users where user_id=$1")).
		WithArgs(ID).
		WillReturnError(NotFound)
	_, err := s.repo.FindByID(ID)
	assert.EqualError(s.T(), NotFound, err.Error())

	u := models.TestUser(s.T())
	u.ID = int(ID)

	assert.NoError(s.T(), u.Encrypt())

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, nickname, avatar from users where user_id=$1")).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(strconv.Itoa(u.ID), u.Nickname, u.Avatar))
	var gotten *models.User
	gotten, err = s.repo.FindByID(ID)

	assert.NotNil(s.T(), gotten)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gotten.ID, u.ID)
	assert.Equal(s.T(), gotten.Nickname, u.Nickname)
	assert.Equal(s.T(), gotten.Avatar, u.Avatar)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}
