package sqlstore

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"patreon/internal/models"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SuiteCreatorRepository struct {
	Suite
	store *Store
}

func (s *SuiteCreatorRepository) SetupSuite() {
	s.InitBD()
	s.store = New(s.DB)
}

func (s *SuiteCreatorRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *SuiteCreatorRepository) TestCreatorRepository_Create() {
	cr := models.TestCreator(s.T())

	cr.ID = 1
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+"RETURNING creator_id")).
		WithArgs(cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(cr.ID)))
	err := s.store.Creator().Create(cr)
	assert.NoError(s.T(), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreator() {
	cr := models.TestCreator(s.T())
	cr.ID = 1
	expected := *cr

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1")).
		WithArgs(cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"}).
			AddRow(strconv.Itoa(cr.ID), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname))

	get, err := s.store.Creator().GetCreator(int64(expected.ID))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreators_AllUsersCreators() {
	creators := models.TestCreators(s.T())

	preapareRows := sqlmock.NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"})

	for index, cr := range creators {
		cr.ID = index
		creators[index] = cr
		preapareRows.AddRow(strconv.Itoa(cr.ID), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname)
	}

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) from creator_profile")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")).
		WillReturnRows(preapareRows)

	get, err := s.store.Creator().GetCreators()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), creators, get)
}

func TestCreatorRepository(t *testing.T) {
	suite.Run(t, new(SuiteCreatorRepository))
}
