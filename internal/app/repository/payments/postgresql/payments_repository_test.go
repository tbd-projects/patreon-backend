package repository_postgresql

import (
	"fmt"
	"patreon/internal/app/models"
	putilits "patreon/internal/app/utilits/postgresql"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"
)

type SuitePaymentsRepository struct {
	models.Suite
	repo *PaymentsRepository
}

func (s *SuitePaymentsRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewPaymentsRepository(s.DB)
}

func (s *SuitePaymentsRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}
func (s *SuitePaymentsRepository) TestPaymentsRepository_GetUserPayments_OK() {
	queryStat := "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
	tableName := "payments"

	payment := models.TestPayment()
	payment.UserID = 0
	creator := models.TestCreator()
	userId := int64(5)

	pag := &models.Pagination{Limit: 10, Offset: 20}
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).
			AddRow(int64(5000)))
	limit, offset, err := putilits.AddPagination(tableName, pag, s.DB)
	query := querySelectUserPayments + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).
			AddRow(int64(5000)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"p.amount", "p.date", "p.creator_id", "u.nickname", "cp.category", "cp.description", "status"}).
			AddRow(payment.Amount, payment.Date, payment.CreatorID, creator.Nickname, creator.Category, creator.Description, payment.Status))
	expRes := []models.UserPayments{
		{
			Payments:           *payment,
			CreatorNickname:    creator.Nickname,
			CreatorCategory:    creator.Category,
			CreatorDescription: creator.Description,
		},
	}
	res, err := s.repo.GetUserPayments(userId, pag)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expRes[0].Payments, res[0].Payments)

}

func (s *SuitePaymentsRepository) TestPaymentsRepository_GetCreatorPayments_OK() {
	queryStat := "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
	tableName := "payments"

	payment := models.TestPayment()
	payment.CreatorID = 0
	user := models.TestUser()
	creatorId := int64(5)

	pag := &models.Pagination{Limit: 10, Offset: 20}
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).
			AddRow(int64(5000)))
	limit, offset, err := putilits.AddPagination(tableName, pag, s.DB)
	query := querySelectCreatorPayments + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).
			AddRow(int64(5000)))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"p.amount", "p.date", "p.users_id", "u.nickname", "status"}).
			AddRow(payment.Amount, payment.Date, payment.UserID, user.Nickname, payment.Status))
	expRes := []models.CreatorPayments{
		{
			Payments:     *payment,
			UserNickname: user.Nickname,
		},
	}
	res, err := s.repo.GetCreatorPayments(creatorId, pag)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), expRes[0].Payments, res[0].Payments)

}
func TestPaymentsRepository(t *testing.T) {
	suite.Run(t, new(SuitePaymentsRepository))
}
