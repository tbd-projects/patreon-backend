package repository_subscribers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"
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

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	awardsName := "daaa"
	price := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectCommit()

	err := s.repo.Create(subscriber, awardsName)
	assert.NoError(s.T(), err)
}


func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_FirstQueryDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"

	awardsName := "daaa"
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).WillReturnError(repository.DefaultErrDB)
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_BeginDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"

	awardsName := "daaa"
	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin().WillReturnError(repository.DefaultErrDB)
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_SecondQueryDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	awardsName := "daaa"
	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnError(repository.DefaultErrDB)
	s.Mock.ExpectRollback()
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_ThirdQueryDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"

	awardsName := "daaa"
	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnError(repository.DefaultErrDB)
	s.Mock.ExpectRollback()
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_CommitDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"

	awardsName := "daaa"
	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectCommit().WillReturnError(repository.DefaultErrDB)
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_CloseRowDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	awardsName := "daaa"
	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))
	s.Mock.ExpectRollback()
	err := s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.CreatorID, awardsName).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))
	s.Mock.ExpectRollback()
	err = s.repo.Create(subscriber, awardsName)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
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

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_RowError_Few() {
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
			AddRow(strconv.Itoa(int(mockRes[2]))).RowError(0, repository.DefaultErrDB))

	_, err := s.repo.GetSubscribers(int64(crId))
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_DBErrorScab_Few() {
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
			AddRow(queryCnt))

	_, err := s.repo.GetSubscribers(int64(crId))
	assert.Error(s.T(), err)
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

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_DBError_Few() {
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
			AddRow(strconv.Itoa(int(mockRes[2]))).RowError(1, repository.DefaultErrDB))

	_, err := s.repo.GetCreators(int64(uId))
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_DBErrorScab_Few() {
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
			AddRow(queryCnt))

	_, err := s.repo.GetCreators(int64(uId))
	assert.Error(s.T(), err)
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

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_Oke() {
	uId := int64(1)
	mockRes := int64(1)
	query := "SELECT count(*) from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes))

	res, err := s.repo.Get(uId, uId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_ZeroSub() {
	uId := int64(1)
	mockRes := int64(0)
	query := "SELECT count(*) from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes))

	res, err := s.repo.Get(uId, uId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), false, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_DBError() {
	uId := int64(1)
	query := "SELECT count(*) from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(repository.DefaultErrDB)

	res, err := s.repo.Get(uId, uId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
	assert.Equal(s.T(), false, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_Oke() {
	uId := int64(1)
	mockRes := int64(1)
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes))

	err := s.repo.Delete(&models.Subscriber{UserID: uId, CreatorID: uId})
	assert.NoError(s.T(), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_DbError() {
	uId := int64(1)

	query := "DELETE from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(repository.DefaultErrDB)

	err := s.repo.Delete(&models.Subscriber{UserID: uId, CreatorID: uId})
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}


func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_ErrorClose() {
	uId := int64(1)
	mockRes := int64(1)
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(uId, uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes).CloseError(repository.DefaultErrDB))

	err := s.repo.Delete(&models.Subscriber{UserID: uId, CreatorID: uId})
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func TestSubscribersRepository(t *testing.T) {
	suite.Run(t, new(SuiteSubscribersRepository))
}
