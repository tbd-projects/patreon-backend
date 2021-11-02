package repository_subscribers

import (
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/lib/pq"

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
	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id, awards_id) VALUES ($1, $2, $3)"

	price := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()
	s.Mock.ExpectCommit()

	err := s.repo.Create(subscriber)
	assert.NoError(s.T(), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_FirstQueryDbError() {
	subscriber := models.TestSubscriber()
	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).WillReturnError(repository.DefaultErrDB)
	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_BeginDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))
	s.Mock.ExpectBegin().WillReturnError(repository.DefaultErrDB)

	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_SecondQueryDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnError(repository.DefaultErrDB)

	s.Mock.ExpectRollback()

	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_ThirdQueryDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id, awards_id) VALUES ($1, $2, $3)"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{})).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnError(repository.DefaultErrDB)

	s.Mock.ExpectRollback()
	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_CommitDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id, awards_id) VALUES ($1, $2, $3)"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectCommit().WillReturnError(repository.DefaultErrDB)
	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_FirstQueryCloseRowDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price).
			CloseError(repository.DefaultErrDB))

	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_SecondQueryCloseRowDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			CloseError(repository.DefaultErrDB))
	s.Mock.ExpectRollback()

	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_Create_ThirdQueryCloseRowDbError() {
	subscriber := models.TestSubscriber()

	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id, awards_id) VALUES ($1, $2, $3)"

	price := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAwardPrice)).
		WithArgs(subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price)).
		RowsWillBeClosed()

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddPayment)).
		WithArgs(price, subscriber.CreatorID, subscriber.UserID).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddSubscribe)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			CloseError(repository.DefaultErrDB))

	s.Mock.ExpectRollback()

	err := s.repo.Create(subscriber)
	assert.Error(s.T(), repository.NewDBError(repository.DefaultErrDB), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_One() {
	crId := 1
	u := models.TestUser()
	u.Login = ""
	u.Password = ""
	mockRes := []models.User{*u}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.users_id", "nickname", "avatar"}).
			AddRow(mockRes[0].ID, mockRes[0].Nickname, mockRes[0].Avatar)).RowsWillBeClosed()

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_Few() {
	crId := 1
	mockRes := []models.User{*models.TestUser(), *models.TestUser(), *models.TestUser()}
	for i := 0; i < len(mockRes); i++ {
		mockRes[i].Login = ""
		mockRes[i].Password = ""
	}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.users_id", "nickname", "avatar"}).
			AddRow(mockRes[0].ID, mockRes[0].Nickname, mockRes[0].Avatar).
			AddRow(mockRes[1].ID, mockRes[1].Nickname, mockRes[1].Avatar).
			AddRow(mockRes[2].ID, mockRes[2].Nickname, mockRes[2].Avatar)).
		RowsWillBeClosed()

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_RowError_Few() {
	crId := 1
	mockRes := []models.User{*models.TestUser(), *models.TestUser(), *models.TestUser()}
	for i := 0; i < len(mockRes); i++ {
		mockRes[i].Login = ""
		mockRes[i].Password = ""
	}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.users_id", "nickname", "avatar"}).
			AddRow(mockRes[0].ID, mockRes[0].Nickname, mockRes[0].Avatar).
			AddRow(mockRes[1].ID, mockRes[1].Nickname, mockRes[1].Avatar).
			AddRow(mockRes[2].ID, mockRes[2].Nickname, mockRes[2].Avatar).RowError(0, repository.DefaultErrDB))

	_, err := s.repo.GetSubscribers(int64(crId))
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_DBErrorScan_Few() {
	crId := 1
	mockRes := []models.User{*models.TestUser(), *models.TestUser(), *models.TestUser()}
	for i := 0; i < len(mockRes); i++ {
		mockRes[i].Login = ""
		mockRes[i].Password = ""
	}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.users_id", "nickname", "avatar"}).
			AddRow("id", "la", "o"))

	_, err := s.repo.GetSubscribers(int64(crId))
	assert.Error(s.T(), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_Ok_NoOne() {
	crId := 1
	mockRes := []models.User{}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{}))

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_CountQueryError() {
	crId := 1
	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	expError := repository.NewDBError(models.BDError)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(models.BDError)

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.Equal(s.T(), expError, err)
	assert.Nil(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetSubscribers_SelectQueryError() {
	crId := 1
	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`
	expError := repository.NewDBError(models.BDError)

	mockRes := []models.User{*models.TestUser()}
	mockRes[0].Login = ""
	mockRes[0].Password = ""

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(crId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(crId).
		WillReturnError(models.BDError)

	res, err := s.repo.GetSubscribers(int64(crId))
	assert.Equal(s.T(), expError, err)
	assert.Nil(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_One() {
	uId := int64(1)
	mockRes := []models.CreatorSubscribe{*models.TestCreatorSubscriber()}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).
			AddRow(len(mockRes))).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.creator_id", "awards_id", "category",
			"description", "nickname", "cp.avatar", "cover"}).
			AddRow(mockRes[0].ID, mockRes[0].AwardsId, mockRes[0].Category, mockRes[0].Description,
				mockRes[0].Nickname, mockRes[0].Avatar, mockRes[0].Cover)).
		RowsWillBeClosed()

	res, err := s.repo.GetCreators(uId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_Few() {
	uId := int64(1)
	mockRes := []models.CreatorSubscribe{*models.TestCreatorSubscriber(),
		*models.TestCreatorSubscriber(), *models.TestCreatorSubscriber()}

	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.creator_id", "awards_id", "category",
			"description", "nickname", "cp.avatar", "cover"}).
			AddRow(mockRes[0].ID, mockRes[0].AwardsId, mockRes[0].Category, mockRes[0].Description,
				mockRes[0].Nickname, mockRes[0].Avatar, mockRes[0].Cover).
			AddRow(mockRes[1].ID, mockRes[1].AwardsId, mockRes[1].Category, mockRes[1].Description,
				mockRes[1].Nickname, mockRes[1].Avatar, mockRes[1].Cover).
			AddRow(mockRes[2].ID, mockRes[2].AwardsId, mockRes[2].Category, mockRes[2].Description,
				mockRes[2].Nickname, mockRes[2].Avatar, mockRes[2].Cover)).
		RowsWillBeClosed()
	res, err := s.repo.GetCreators(uId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_DBError_Few() {
	uId := int64(1)
	mockRes := []models.CreatorSubscribe{*models.TestCreatorSubscriber(),
		*models.TestCreatorSubscriber(), *models.TestCreatorSubscriber()}
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.creator_id", "awards_id", "category",
			"description", "nickname", "cp.avatar", "cover"}).
			AddRow(mockRes[0].ID, mockRes[0].AwardsId, mockRes[0].Category, mockRes[0].Description,
				mockRes[0].Nickname, mockRes[0].Avatar, mockRes[0].Cover).
			AddRow(mockRes[1].ID, mockRes[1].AwardsId, mockRes[1].Category, mockRes[1].Description,
				mockRes[1].Nickname, mockRes[1].Avatar, mockRes[1].Cover).
			AddRow(mockRes[2].ID, mockRes[2].AwardsId, mockRes[2].Category, mockRes[2].Description,
				mockRes[2].Nickname, mockRes[2].Avatar, mockRes[2].Cover).RowError(2, sqlErr)).
		RowsWillBeClosed()
	res, err := s.repo.GetCreators(uId)

	assert.Nil(s.T(), res)
	assert.Error(s.T(), expErr, err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_DBErrorScan_Few() {
	uId := int64(1)
	mockRes := []models.CreatorSubscribe{*models.TestCreatorSubscriber(),
		*models.TestCreatorSubscriber(), *models.TestCreatorSubscriber()}

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)

	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).
			AddRow(len(mockRes))).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.creator_id", "awards_id", "category",
			"description", "nickname", "cp.avatar", "cover"}).
			AddRow(mockRes[0].ID, mockRes[0].AwardsId, mockRes[0].Category, mockRes[0].Description,
				mockRes[0].Nickname, mockRes[0].Avatar, mockRes[0].Cover).
			AddRow(mockRes[1].ID, mockRes[1].AwardsId, mockRes[1].Category, mockRes[1].Description,
				mockRes[1].Nickname, mockRes[1].Avatar, mockRes[1].Cover).RowError(1, sqlErr).
			AddRow(mockRes[2].ID, mockRes[2].AwardsId, mockRes[2].Category, mockRes[2].Description,
				mockRes[2].Nickname, mockRes[2].Avatar, mockRes[2].Cover))

	res, err := s.repo.GetCreators(uId)
	assert.Equal(s.T(), expErr, err)
	assert.Nil(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_CountQueryError() {
	uId := int64(1)
	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(sqlErr)

	res, err := s.repo.GetCreators(uId)
	assert.Equal(s.T(), expErr, err)
	assert.Nil(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_SelectQueryError() {
	uId := int64(1)

	mockRes := []models.CreatorSubscribe{*models.TestCreatorSubscriber()}
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)

	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(len(mockRes)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(sqlErr)

	res, err := s.repo.GetCreators(uId)
	assert.Equal(s.T(), expErr, err)
	assert.Nil(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_GetCreators_Ok_NoOne() {
	uId := int64(1)
	var mockRes []models.CreatorSubscribe

	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).
			AddRow(len(mockRes))).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(uId).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"s.creator_id", "awards_id", "category",
			"description", "nickname", "cp.avatar", "cover"})).
		RowsWillBeClosed()

	res, err := s.repo.GetCreators(uId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_Oke() {
	mockRes := int64(1)
	subscriber := models.TestSubscriber()
	query := "SELECT count(*) as cnt from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes))

	res, err := s.repo.Get(subscriber)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_ZeroSub() {
	mockRes := int64(0)
	subscriber := models.TestSubscriber()
	query := "SELECT count(*) as cnt from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes))

	res, err := s.repo.Get(subscriber)
	assert.NoError(s.T(), err)
	assert.False(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Get_DBError() {
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	subscriber := models.TestSubscriber()
	query := "SELECT count(*) as cnt from subscribers where users_id = $1 and creator_id = $2"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID).
		WillReturnError(sqlErr)

	res, err := s.repo.Get(subscriber)
	assert.Equal(s.T(), expErr, err)
	assert.False(s.T(), res)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_Oke() {
	subscriber := models.TestSubscriber()
	mockRes := int64(1)
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2 and awards_id = $3"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes)).
		RowsWillBeClosed()

	err := s.repo.Delete(subscriber)
	assert.NoError(s.T(), err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_DbError() {
	subscriber := models.TestSubscriber()
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2 and awards_id = $3"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnError(sqlErr)

	err := s.repo.Delete(subscriber)
	assert.Equal(s.T(), expErr, err)
}

func (s *SuiteSubscribersRepository) TestSubscribersRepository_Delete_ErrorClose() {
	subscriber := models.TestSubscriber()
	mockRes := int64(1)
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2 and awards_id = $3"

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(subscriber.UserID, subscriber.CreatorID, subscriber.AwardID).
		WillReturnError(nil).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(mockRes).CloseError(sqlErr)).
		RowsWillBeClosed()

	err := s.repo.Delete(subscriber)
	assert.Equal(s.T(), expErr, err)
}

func TestSubscribersRepository(t *testing.T) {
	suite.Run(t, new(SuiteSubscribersRepository))
}
