package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/lib/pq"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SuiteUserRepository struct {
	models.Suite
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
	u := models.TestUser()

	u.ID = 1
	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) "+
		"RETURNING users_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(u.ID))))

	err := s.repo.Create(u)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) "+
		"RETURNING users_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnError(models.BDError)

	err = s.repo.Create(u)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) "+
		"RETURNING users_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnError(&pq.Error{Code: codeDuplicateVal, Constraint: loginConstraint})

	err = s.repo.Create(u)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), LoginAlreadyExist, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) "+
		"RETURNING users_id")).
		WithArgs(u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).
		WillReturnError(&pq.Error{Code: codeDuplicateVal, Constraint: nicknameConstraint})

	err = s.repo.Create(u)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), NicknameAlreadyExist, err)
}

func (s *SuiteUserRepository) TestUserRepository_FindByLogin() {
	login := "mail1999"
	query := "SELECT users_id, login, nickname, avatar, encrypted_password " +
		"from users where login=$1"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(login).
		WillReturnError(models.BDError)
	_, err := s.repo.FindByLogin(login)
	assert.EqualError(s.T(), repository.NewDBError(models.BDError), err.Error())

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(login).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.FindByLogin(login)
	assert.EqualError(s.T(), repository.NotFound, err.Error())

	u := models.TestUser()
	u.Login = login

	assert.NoError(s.T(), u.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(login).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "nickname",
			"avatar", "encrypted_password"}).
			AddRow(strconv.Itoa(int(u.ID)), u.Login, u.Nickname, u.Avatar,
				u.EncryptedPassword))
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
	query := `SELECT users_id, login, nickname, users.avatar, encrypted_password, cp.creator_id IS NOT NULL
	from users LEFT JOIN creator_profile AS cp ON (users.users_id = cp.creator_id) where users_id=$1`
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(ID).
		WillReturnError(sql.ErrNoRows)
	_, err := s.repo.FindByID(ID)
	assert.EqualError(s.T(), repository.NotFound, err.Error())

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(ID).
		WillReturnError(models.BDError)
	_, err = s.repo.FindByID(ID)
	assert.EqualError(s.T(), repository.NewDBError(models.BDError), err.Error())

	u := models.TestUser()
	u.ID = ID

	assert.NoError(s.T(), u.Encrypt())

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "nickname",
			"avatar", "encrypted_password", "creator_id"}).
			AddRow(strconv.Itoa(int(u.ID)), u.Login, u.Nickname, u.Avatar,
				u.EncryptedPassword, true))
	var gotten *models.User
	gotten, err = s.repo.FindByID(ID)

	assert.NotNil(s.T(), gotten)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), gotten.ID, u.ID)
	assert.Equal(s.T(), gotten.Nickname, u.Nickname)
	assert.Equal(s.T(), gotten.Avatar, u.Avatar)
	assert.True(s.T(), gotten.HaveCreator)
}
func (s *SuiteUserRepository) TestUserRepository_UpdateAvatar_Correct() {
	user := models.TestUser()
	newAvatar := "newAvatar.png"
	s.Mock.ExpectQuery(regexp.QuoteMeta("UPDATE users SET avatar = $1 WHERE users_id = $2")).
		WithArgs(newAvatar, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	res := s.repo.UpdateAvatar(user.ID, newAvatar)
	assert.NoError(s.T(), res)
}
func (s *SuiteUserRepository) TestUserRepository_UpdateAvatar_Error() {
	user := models.TestUser()
	newAvatar := "newAvatar.png"
	s.Mock.ExpectQuery(regexp.QuoteMeta("UPDATE users SET avatar = $1 WHERE users_id = $2")).
		WithArgs(newAvatar, user.ID).
		WillReturnError(models.BDError)
	res := s.repo.UpdateAvatar(user.ID, newAvatar)
	assert.Error(s.T(), models.BDError, res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdatePassword_Correct() {
	user := models.TestUser()
	assert.NoError(s.T(), user.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta("UPDATE users SET encrypted_password = $1 WHERE users_id = $2")).
		WithArgs(user.EncryptedPassword, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	res := s.repo.UpdatePassword(user.ID, user.EncryptedPassword)
	assert.NoError(s.T(), res)
}
func (s *SuiteUserRepository) TestUserRepository_UpdatePassword_Error() {
	user := models.TestUser()
	assert.NoError(s.T(), user.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta("UPDATE users SET encrypted_password = $1 WHERE users_id = $2")).
		WithArgs(user.EncryptedPassword, user.ID).
		WillReturnError(models.BDError)

	res := s.repo.UpdatePassword(user.ID, user.EncryptedPassword)
	assert.Error(s.T(), models.BDError, res)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}
