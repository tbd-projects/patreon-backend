package repository_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"github.com/lib/pq"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
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
	runFunc := func(input ...interface{}) (res []interface{}) {
		oldNick, _ := input[0].(*models.User)
		err := s.repo.Create(oldNick)
		return []interface{}{err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(u.ID),
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Query:   createSettingsQuery,
					Err:     nil,
					Args:    []driver.Value{u.ID},
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "ErrBegin",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     repository.DefaultErrDB,
					RunType: models.TransBegin,
				},
			},
		},
		{
			Name: "ErrCreate",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query:   createQuery,
					Err:     repository.DefaultErrDB,
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrCreateDupleLogin",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     LoginAlreadyExist,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err: &pq.Error{
						Code:       codeDuplicateVal,
						Constraint: loginConstraint,
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrCreateDupleNickname",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     NicknameAlreadyExist,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err: &pq.Error{
						Code:       codeDuplicateVal,
						Constraint: nicknameConstraint,
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrCreateOtherError",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err: &pq.Error{
						Code:       "sadasd",
						Constraint: nicknameConstraint,
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrCreateSettings",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(u.ID),
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Query:   createSettingsQuery,
					Err:     repository.DefaultErrDB,
					Args:    []driver.Value{u.ID},
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrCommit",
			Args: []interface{}{u},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query: createQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(u.ID),
					},
					Args:    []driver.Value{u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage},
					RunType: models.Query,
				},
				{
					Query:   createSettingsQuery,
					Err:     nil,
					Args:    []driver.Value{u.ID},
					RunType: models.Exec,
				},
				{
					Err:     repository.DefaultErrDB,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
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
	query := `UPDATE users SET avatar = $1 WHERE users_id = $2`
	user := models.TestUser()
	newAvatar := "newAvatar.png"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(newAvatar, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	res := s.repo.UpdateAvatar(user.ID, newAvatar)
	assert.NoError(s.T(), res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdateAvatar_CloseError() {
	query := `UPDATE users SET avatar = $1 WHERE users_id = $2`
	user := models.TestUser()
	newAvatar := "newAvatar.png"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(newAvatar, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))
	res := s.repo.UpdateAvatar(user.ID, newAvatar)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdateAvatar_Error() {
	query := `UPDATE users SET avatar = $1 WHERE users_id = $2`
	user := models.TestUser()
	newAvatar := "newAvatar.png"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(newAvatar, user.ID).
		WillReturnError(models.BDError)
	res := s.repo.UpdateAvatar(user.ID, newAvatar)
	assert.Error(s.T(), models.BDError, res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdatePassword_Correct() {
	query := `UPDATE users SET encrypted_password = $1 WHERE users_id = $2`
	user := models.TestUser()
	assert.NoError(s.T(), user.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(user.EncryptedPassword, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	res := s.repo.UpdatePassword(user.ID, user.EncryptedPassword)
	assert.NoError(s.T(), res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdatePassword_CloseError() {
	query := `UPDATE users SET encrypted_password = $1 WHERE users_id = $2`
	user := models.TestUser()
	assert.NoError(s.T(), user.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(user.EncryptedPassword, user.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))

	res := s.repo.UpdatePassword(user.ID, user.EncryptedPassword)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdatePassword_Error() {
	query := `UPDATE users SET encrypted_password = $1 WHERE users_id = $2`
	user := models.TestUser()
	assert.NoError(s.T(), user.Encrypt())
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(user.EncryptedPassword, user.ID).
		WillReturnError(models.BDError)

	res := s.repo.UpdatePassword(user.ID, user.EncryptedPassword)
	assert.Error(s.T(), models.BDError, res)
}

func (s *SuiteUserRepository) TestUserRepository_UpdateNickname() {
	oldNickName := "d"
	newNickName := "32"
	runFunc := func(input ...interface{}) (res []interface{}) {
		oldNick, _ := input[0].(string)
		newNick, _ := input[1].(string)
		err := s.repo.UpdateNickname(oldNick, newNick)
		return []interface{}{err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{oldNickName, newNickName},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: updateNicknameQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnResult: driver.RowsAffected(1),
					},
					Args:    []driver.Value{newNickName, oldNickName},
					RunType: models.Exec,
				},
			},
		},
		{
			Name: "NothingAffected",
			Args: []interface{}{oldNickName, newNickName},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: updateNicknameQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnResult: driver.RowsAffected(0),
					},
					Args:    []driver.Value{newNickName, oldNickName},
					RunType: models.Exec,
				},
			},
		},
		{
			Name: "BDError",
			Args: []interface{}{oldNickName, newNickName},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(models.BDError),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   updateNicknameQuery,
					Err:     models.BDError,
					Rows:    nil,
					Args:    []driver.Value{newNickName, oldNickName},
					RunType: models.Exec,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteUserRepository) TestUserRepository_IsAllowedAward() {
	userId := int64(1)
	awardId := int64(2)
	runFunc := func(input ...interface{}) (res []interface{}) {
		userId, _ := input[0].(int64)
		awardId, _ := input[1].(int64)
		is, err := s.repo.IsAllowedAward(userId, awardId)
		return []interface{}{is, err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{userId, awardId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{true},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: isAllowedAwardQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
					},
					Args:    []driver.Value{awardId, userId},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{userId, awardId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   isAllowedAwardQuery,
					Err:     sql.ErrNoRows,
					Rows:    nil,
					Args:    []driver.Value{awardId, userId},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BDError",
			Args: []interface{}{userId, awardId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(models.BDError),
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   isAllowedAwardQuery,
					Err:     models.BDError,
					Rows:    nil,
					Args:    []driver.Value{awardId, userId},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "CorrectNotAllowed",
			Args: []interface{}{userId, awardId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: isAllowedAwardQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"count"}).AddRow(0),
					},
					Args:    []driver.Value{awardId, userId},
					RunType: models.Query,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteUserRepository) TestUserRepository_FindByNickname() {
	user := models.TestUser()
	user.Password = ""
	nickname := "dor"
	runFunc := func(input ...interface{}) (res []interface{}) {
		nickname, _ := input[0].(string)
		user, err := s.repo.FindByNickname(nickname)
		return []interface{}{user, err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{nickname},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{user},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: findByNicknameQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id", "login", "nickname", "avatar", "pass", "haveCreator"}).
							AddRow(user.ID, user.Login, user.Nickname, user.Avatar,
								user.EncryptedPassword, user.HaveCreator),
					},
					Args:    []driver.Value{nickname},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{nickname},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{(*models.User)(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   findByNicknameQuery,
					Err:     sql.ErrNoRows,
					Rows:    nil,
					Args:    []driver.Value{nickname},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BdError",
			Args: []interface{}{nickname},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{(*models.User)(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   findByNicknameQuery,
					Err:     repository.DefaultErrDB,
					Rows:    nil,
					Args:    []driver.Value{nickname},
					RunType: models.Query,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}
func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}
