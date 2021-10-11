package repository_creator

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type SuiteCreatorRepository struct {
	models.Suite
	repo *CreatorRepository
}

func (s *SuiteCreatorRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewCreatorRepository(s.DB)
}

func (s *SuiteCreatorRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteCreatorRepository) TestCreatorRepository_Create() {
	cr := models.TestCreator()

	cr.ID = 1
	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+"RETURNING creator_id")).
		WithArgs(cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(cr.ID))))
	id, err := s.repo.Create(cr)
	assert.Equal(s.T(), id, cr.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+"RETURNING creator_id")).
		WithArgs(cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreator() {
	cr := models.TestCreator()
	cr.ID = 1
	expected := *cr

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1")).
		WithArgs(cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname))

	get, err := s.repo.GetCreator(expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1")).
		WithArgs(cr.ID).WillReturnError(sql.ErrNoRows)

	_, err = s.repo.GetCreator(expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NotFound, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1")).
		WithArgs(cr.ID).WillReturnError(models.BDError)

	_, err = s.repo.GetCreator(expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreators_AllUsersCreators() {
	creators := models.TestCreators()

	preapareRows := sqlmock.NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"})

	for index, cr := range creators {
		cr.ID = int64(index)
		creators[index] = cr
		preapareRows.AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname)
	}

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) from creator_profile")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")).
		WillReturnRows(preapareRows)

	get, err := s.repo.GetCreators()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), creators, get)

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) from creator_profile")).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) from creator_profile")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func TestCreatorRepository(t *testing.T) {
	suite.Run(t, new(SuiteCreatorRepository))
}
