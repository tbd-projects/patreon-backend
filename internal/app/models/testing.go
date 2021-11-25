package models

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

var BDError = errors.New("BD error")

type typeQuery int64

const (
	Exec = iota
	Query
	TransBegin
	TransCommit
	TransRollback
)

const (
	TransactionContext = "trans"
)

var (
	//NotValidTestCase = errors.New("Not found transaction in context queries ")
	UnknownQueryType = errors.New("Unknown query type, only can be: Exec, Query, TransBegin, TransRollback ")
)

type TestRowError struct {
	Err error
	Row int
}

type TestRow struct {
	ReturnRows   *sqlmock.Rows
	ReturnResult driver.Result
	RowError     *TestRowError
	RowClose     error
}

func (row *TestRow) addRow(expect *sqlmock.ExpectedQuery) {
	if row.RowError == nil && row.RowClose == nil {
		expect.WillReturnRows(row.ReturnRows)
		return
	}

	if row.RowClose != nil {
		expect.WillReturnRows(row.ReturnRows.CloseError(row.RowClose))
		return
	}

	expect.WillReturnRows(row.ReturnRows.RowError(row.RowError.Row, row.RowError.Err))
}

type TestQuery struct {
	Query   string
	Err     error
	Rows    *TestRow
	Args    []driver.Value
	RunType typeQuery
}

type TestExpected struct {
	CheckError      bool // mean only check error without compare
	HaveError       bool // in returning values end of values is error
	ExpectedErr     error
	ExpectedReturns []interface{}
}

type TestCase struct {
	Name     string
	Args     []interface{}
	Queries  []TestQuery
	RunFunc  func(...interface{}) []interface{}
	Expected TestExpected
}

type Suite struct {
	suite.Suite
	DB   *sqlx.DB
	Mock sqlmock.Sqlmock
}

func (s *Suite) transBeginQuery(_ context.Context, returnErr error) error {
	s.Mock.ExpectBegin().WillReturnError(returnErr)
	return nil
}

func (s *Suite) transRollbackQuery(_ context.Context, returnErr error) error {
	s.Mock.ExpectRollback().WillReturnError(returnErr)
	return nil
}

func (s *Suite) transCommitQuery(_ context.Context, returnErr error) error {
	s.Mock.ExpectCommit().WillReturnError(returnErr)
	return nil
}

func (s *Suite) execQuery(_ context.Context, query TestQuery) error {
	exec := s.Mock.ExpectExec(regexp.QuoteMeta(query.Query)).WithArgs(query.Args...)

	if query.Err != nil {
		exec.WillReturnError(query.Err)
	} else {
		if query.Rows == nil || query.Rows.ReturnResult == nil {
			exec.WillReturnResult(driver.ResultNoRows)
		} else {
			exec.WillReturnResult(query.Rows.ReturnResult)
		}
	}
	return nil
}

func (s *Suite) baseQuery(_ context.Context, query TestQuery) error {
	exec := s.Mock.ExpectQuery(regexp.QuoteMeta(query.Query)).WithArgs(query.Args...)

	if query.Err != nil {
		exec.WillReturnError(query.Err)
	} else {
		query.Rows.addRow(exec)
	}
	return nil
}

func (s *Suite) runQuery(ctx context.Context, query TestQuery) error {
	switch query.RunType {
	case Exec:
		return s.execQuery(ctx, query)
	case Query:
		return s.baseQuery(ctx, query)
	case TransCommit:
		return s.transCommitQuery(ctx, query.Err)
	case TransRollback:
		return s.transRollbackQuery(ctx, query.Err)
	case TransBegin:
		return s.transBeginQuery(ctx, query.Err)
	default:
		return UnknownQueryType
	}
}

func (s *Suite) checkExpected(res []interface{}, expected TestExpected, caseName string) {
	if expected.HaveError {
		size := len(res)
		require.NotZerof(s.T(), size, "Testcase with name: %s, return nothing, but wait return error", caseName)
		gottedError, ok := res[size-1].(error)
		if !ok && gottedError != nil {
			require.Failf(s.T(), "Last value not error, but expected error",
				"Testcase with name: %s", caseName)
		}

		if expected.CheckError {
			assert.Error(s.T(), gottedError, "Testcase with name: %s", caseName)
		} else {
			if expected.ExpectedErr == nil {
				assert.NoError(s.T(), gottedError,
					"Testcase with name: %s", caseName)
			} else {
				assert.EqualError(s.T(), gottedError, expected.ExpectedErr.Error(),
					"Testcase with name: %s", caseName)
			}
		}
		res = res[:size-1]
	}

	require.Equalf(s.T(), len(res), len(expected.ExpectedReturns),
		"Testcase with name: %s, different len of expected and gotten return values", caseName)

	for i, expected := range expected.ExpectedReturns {
		assert.Equalf(s.T(), expected, res[i], "Testcase with name: %s", caseName)
	}
}

func (s *Suite) RunTestCase(test TestCase) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			assert.Failf(t, "Error testcase: ", "%s %s", test.Name, r.(error).Error())
		}
	}(s.T())

	ctx := context.Background()
	for _, query := range test.Queries {
		err := s.runQuery(ctx, query)
		require.NoErrorf(s.T(), err, "Testcase with name: %s", test.Name)
	}
	res := test.RunFunc(test.Args...)

	s.checkExpected(res, test.Expected, test.Name)

}

func (s *Suite) InitBD() {
	s.T().Helper()

	var err error
	s.DB, s.Mock, err = sqlmock.Newx()
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestUser() *User {
	return &User{
		Login:    "student1999",
		Password: "1!2!3!",
		Nickname: "patron",
		Avatar:   "default.png",
	}

}
func TestUsers() []User {
	u1 := User{
		Login:    "test1",
		Nickname: "test1",
		Password: "123456",
	}
	u2 := User{
		Login:    "test2",
		Nickname: "test2",
		Password: "123456",
	}
	u3 := User{
		Login:    "test3",
		Nickname: "test3",
		Password: "123456",
	}
	return []User{u1, u2, u3}
}

func TestCreator() *Creator {
	return &Creator{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
	}
}

func TestCreatorWithAwards() *CreatorWithAwards {
	return &CreatorWithAwards{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
		AwardsId:    1,
	}
}

func TestCreatorSubscriber() *CreatorSubscribe {
	return &CreatorSubscribe{
		ID:          1,
		Nickname:    "doggy2005",
		Category:    "podcasts",
		Description: "i love podcasts",
		AwardsId:    1,
	}
}

func TestCreators() []Creator {
	cr1 := Creator{
		ID:          1,
		Category:    "podcasts",
		Description: "i love podcasts",
	}
	cr2 := Creator{
		ID:          2,
		Category:    "blog",
		Description: "i love podcasts",
	}
	cr3 := Creator{
		ID:          3,
		Category:    "movies",
		Description: "i love podcasts",
	}
	return []Creator{cr1, cr2, cr3}
}
func TestSubscriber() *Subscriber {
	return &Subscriber{
		ID:        1,
		UserID:    1,
		CreatorID: 2,
		AwardID:   3,
	}

}
func TestCreatorSubscribe() *CreatorSubscribe {
	return &CreatorSubscribe{
		ID:          1,
		Category:    "cat",
		Nickname:    "nick",
		Description: "desc",
		AwardsId:    2,
	}

}
func TestAward() *Award {
	return &Award{
		ID:          1,
		Name:        "award",
		Description: "description",
		Price:       100,
		CreatorId:   1,
		Cover:       "not found",
	}
}
func TestLike() *Like {
	return &Like{
		ID:     1,
		Value:  1,
		PostId: 2,
		UserId: 3,
	}
}
func TestUpdatePost() *UpdatePost {
	return &UpdatePost{
		ID:          1,
		Title:       "Title",
		Awards:      1,
		Description: "jfnagd",
	}
}
func TestCreatePost() *CreatePost {
	return &CreatePost{
		ID:          1,
		Title:       "Title",
		Awards:      1,
		Description: "jfnagd",
		CreatorId:   1,
	}
}
func TestAttachWithoutLevel() *AttachWithoutLevel {
	return &AttachWithoutLevel{
		ID:     1,
		PostId: 1,
		Value:  "jfnagd",
		Type:   Image,
	}
}

func TestPayment() *Payments {
	return &Payments{
		Amount:    100,
		Date:      time.Now(),
		CreatorID: 1,
		UserID:    11,
	}
}

func TestAttach() *Attach {
	return &Attach{
		Id:    1,
		Level: 1,
		Value: "jfnagd",
		Type:  Image,
	}
}
