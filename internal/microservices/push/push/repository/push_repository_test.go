package repository

import (
	"database/sql"
	"database/sql/driver"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"testing"
)

type SuitePushRepository struct {
	models.Suite
	repo *PushRepository
	data models.AttachWithoutLevel
}

func (s *SuitePushRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewPushRepository(s.DB)
	s.data.Value = "asd"
	s.data.Type = "image"
	s.data.ID = 12
	s.data.PostId = 2
}

func (s *SuitePushRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuitePushRepository) TestPushRepository_GetAwardsNameAndPrice() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, sec, err := s.repo.GetAwardsNameAndPrice(input[0].(int64))
		return []interface{}{values, sec, err}
	}
	name := "dore"
	price := int64(2)

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{price},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{name, price},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetAwardsNameAndPriceQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"name", "price"}).
							AddRow(name, price),
					},
					RunType: models.Query,
					Args:    []driver.Value{price},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{price},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{"", int64(-1)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetAwardsNameAndPriceQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{price},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{price},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{"", int64(-1)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetAwardsNameAndPriceQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{price},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_GetCreatorPostAndTitle() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, sec, err := s.repo.GetCreatorPostAndTitle(input[0].(int64))
		return []interface{}{values, sec, err}
	}
	title := "dore"
	creatorId := int64(2)

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{creatorId, title},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetCreatorPostAndTitleQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"name", "price"}).
							AddRow(creatorId, title),
					},
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{int64(-1), ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetCreatorPostAndTitleQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{int64(-1), ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetCreatorPostAndTitleQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_CheckCreatorForGetCommentPush() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, err := s.repo.CheckCreatorForGetCommentPush(input[0].(int64))
		return []interface{}{values, err}
	}
	creatorId := int64(2)

	testings := []models.TestCase{
		{
			Name: "CorrectTrue",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{true},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: CheckCreatorForGetCommentPushQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"get"}).
							AddRow(true),
					},
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "CorrectFalse",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: CheckCreatorForGetCommentPushQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"get"}).
							AddRow(false),
					},
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   CheckCreatorForGetCommentPushQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   CheckCreatorForGetCommentPushQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_CheckCreatorForGetSubPush() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, err := s.repo.CheckCreatorForGetSubPush(input[0].(int64))
		return []interface{}{values, err}
	}
	creatorId := int64(2)

	testings := []models.TestCase{
		{
			Name: "CorrectTrue",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{true},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: CheckCreatorForGetSubPushQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"get"}).
							AddRow(true),
					},
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "CorrectFalse",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: CheckCreatorForGetSubPushQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"get"}).
							AddRow(false),
					},
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   CheckCreatorForGetSubPushQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{false},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   CheckCreatorForGetSubPushQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{creatorId},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_GetSubUserForPushPost() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, err := s.repo.GetSubUserForPushPost(input[0].(int64))
		return []interface{}{values, err}
	}
	id := int64(2)

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{[]int64{id}},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetSubUserForPushPostQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"user"}).
							AddRow(id),
					},
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetSubUserForPushPostQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "ErrorScan",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetSubUserForPushPostQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"user"}).
							AddRow("dore"),
					},
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "RowError",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetSubUserForPushPostQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"user"}).
							AddRow(id),
						RowError: &models.TestRowError{
							Row: 0,
							Err: repository.DefaultErrDB,
						},
					},
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_GetUserNameAndAvatar() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, sec, err := s.repo.GetUserNameAndAvatar(input[0].(int64))
		return []interface{}{values, sec, err}
	}
	avatar := "dore"
	username := "dore"
	id := int64(1)

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{username, avatar},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetUserNameAndAvatarQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"name", "price"}).
							AddRow(username, avatar),
					},
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{"", ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetUserNameAndAvatarQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{"", ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetUserNameAndAvatarQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuitePushRepository) TestPushRepository_GetCreatorNameAndAvatar() {
	runFunc := func(input ...interface{}) (res []interface{}) {
		values, sec, err := s.repo.GetCreatorNameAndAvatar(input[0].(int64))
		return []interface{}{values, sec, err}
	}
	avatar := "dore"
	username := "dore"
	id := int64(1)

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{username, avatar},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: GetCreatorNameAndAvatarQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"name", "price"}).
							AddRow(username, avatar),
					},
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "Err",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{"", ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetCreatorNameAndAvatarQuery,
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{id},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{"", ""},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   GetCreatorNameAndAvatarQuery,
					Err:     sql.ErrNoRows,
					RunType: models.Query,
					Args:    []driver.Value{id},
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func TestAttachesRepository(t *testing.T) {
	suite.Run(t, new(SuitePushRepository))
}
