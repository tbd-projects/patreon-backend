package repository_subscribers

import (
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
)

type SuiteSubscribersRepository struct {
	models.Suite
	repo *SubscribersRepository
}

func (s *SuiteSubscribersRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewSubscribersRepository(s.DB)
}

func (s *SuiteSubscribersRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_Ok() {
	subscriber := models.TestSubscriber()

	expQuery := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"

	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuery)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	err := s.repo.Create(subscriber)
	assert.NoError(s.T(), err)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_DbError() {
	subscriber := models.TestSubscriber()

	expQuery := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"
	expError := repository.NewDBError(models.BDError)
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuery)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnError(expError.ExternalErr)
	err := s.repo.Create(subscriber)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expError, err)

}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_One() {
	crId := 1
	mockRes := []int64{1}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	queryGet := "SELECT users_id from subscribers WHERE creator_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"users_id"}).
			AddRow(strconv.Itoa(int(mockRes[0]))))

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_Few() {
	crId := 1
	mockRes := []int64{1, 2, 3}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	queryGet := "SELECT users_id from subscribers WHERE creator_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"users_id"}).
			AddRow(strconv.Itoa(int(mockRes[0]))).
			AddRow(strconv.Itoa(int(mockRes[1]))).
			AddRow(strconv.Itoa(int(mockRes[2]))))

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_NoOne() {
	crID := 1
	mockRes := []int64{}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	queryGet := "SELECT users_id from subscribers WHERE creator_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(crID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(crID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"users_id"}))
	res, err := s.repo.GetSubscribers(int64(crID))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_CountQueryError() {
	crId := 1
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	expError := repository.NewDBError(models.BDError)
	expRes := []int64{}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(crId).
		WillReturnError(models.BDError)

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.Equal(s.T(), expError, err)
	assert.Equal(s.T(), expRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_SelectQueryError() {
	crID := 1
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	expError := repository.NewDBError(models.BDError)
	mockRes := []int64{1}
	queryGet := "SELECT users_id from subscribers WHERE creator_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(crID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(crID).
		WillReturnError(models.BDError)

	res, err := s.repo.GetSubscribers(int64(crID))
	assert.Equal(s.T(), expError, err)
	assert.Equal(s.T(), []int64{}, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_One() {
	uId := 1
	mockRes := []int64{1}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	queryGet := "SELECT creator_id from subscribers WHERE users_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).
			AddRow(strconv.Itoa(int(mockRes[0]))))
	res, err := s.repo.GetCreators(int64(uId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_Few() {
	uId := 1
	mockRes := []int64{1, 2, 3}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	queryGet := "SELECT creator_id from subscribers WHERE users_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).
			AddRow(strconv.Itoa(int(mockRes[0]))).
			AddRow(strconv.Itoa(int(mockRes[1]))).
			AddRow(strconv.Itoa(int(mockRes[2]))))

	res, err := s.repo.GetCreators(int64(uId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_CountQueryError() {
	uId := 1
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	expError := repository.NewDBError(models.BDError)
	expRes := []int64{}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(uId).
		WillReturnError(models.BDError)

	res, err := s.repo.GetCreators(int64(uId))
	assert.Equal(s.T(), expError, err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_SelectQueryError() {
	uId := 1
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	expError := repository.NewDBError(models.BDError)
	mockRes := []int64{1}
	queryGet := "SELECT creator_id from subscribers WHERE users_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(uId).
		WillReturnError(models.BDError)

	res, err := s.repo.GetCreators(int64(uId))
	assert.Equal(s.T(), expError, err)
	assert.Equal(s.T(), []int64{}, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_NoOne() {
	uId := 1
	mockRes := []int64{}
	queryCnt := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	queryGet := "SELECT creator_id from subscribers WHERE users_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCnt)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGet)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}))
	res, err := s.repo.GetCreators(int64(uId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func TestSubscribersRepository(t *testing.T) {
	suite.Run(t, new(SuiteSubscribersRepository))
}
