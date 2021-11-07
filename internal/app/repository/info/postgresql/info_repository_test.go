package repository_postgresql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"
)

type SuiteInfoRepository struct {
	models.Suite
	repo     *InfoRepository
	testData *models.Info
}

func (s *SuiteInfoRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewInfoRepository(s.DB)
	s.testData = &models.Info{TypePostData: []string{"don", "con"}, Category: []string{"gfy", "ton"}}
}

func (s *SuiteInfoRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_Correct() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(s.testData.Category[0]).
			AddRow(s.testData.Category[1]))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryTypeDataGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"type"}).
			AddRow(s.testData.TypePostData[0]).
			AddRow(s.testData.TypePostData[1]))
	res, err := s.repo.Get()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.testData, res)
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_CategoryIncorrectScan() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(sql.NullString{Valid: false}).
			AddRow(s.testData.Category[1]))
	_, err := s.repo.Get()
	assert.Error(s.T(), err)
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_TypeIncorrectScan() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(s.testData.Category[0]).
			AddRow(s.testData.Category[1]))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryTypeDataGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"type"}).
			AddRow(sql.NullString{Valid: false}).
			AddRow(s.testData.TypePostData[1]))
	_, err := s.repo.Get()
	assert.Error(s.T(), err)
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_CategoryRowError() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(s.testData.Category[0]).
			AddRow(s.testData.Category[1]).RowError(0, repository.DefaultErrDB))
	_, err := s.repo.Get()
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_TypeRowError() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(s.testData.Category[0]).
			AddRow(s.testData.Category[1]))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryTypeDataGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"type"}).
			AddRow(s.testData.TypePostData[0]).
			AddRow(s.testData.TypePostData[1]).RowError(0, repository.DefaultErrDB))
	_, err := s.repo.Get()
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_CategoryBdError() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnError(repository.DefaultErrDB)
	_, err := s.repo.Get()
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteInfoRepository) TestInfoRepositoryGet_TypeDbError() {
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryGet)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(s.testData.Category[0]).
			AddRow(s.testData.Category[1]))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryTypeDataGet)).
		WithArgs().
		WillReturnError(repository.DefaultErrDB)
		_, err := s.repo.Get()
		assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func TestInfoRepository(t *testing.T) {
	suite.Run(t, new(SuiteInfoRepository))
}
