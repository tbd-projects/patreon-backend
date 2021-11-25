package repository_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	rp "patreon/internal/app/repository"
	postgresql_utilits "patreon/internal/app/utilits/postgresql"
	"patreon/pkg/utils"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"

	"github.com/stretchr/testify/assert"
)

type SuiteCreatorRepository struct {
	models.Suite
	repo *CreatorRepository
}

type CustomRows sqlmock.Rows

func (t *CustomRows) Copy() *sqlmock.Rows {
	res := sqlmock.Rows(*t)
	return &res
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
	categoryId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryCreate)).
		WithArgs(cr.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(categoryId))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreate)).
		WithArgs(cr.ID, categoryId, cr.Description, app.DefaultImage, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(cr.ID))))
	id, err := s.repo.Create(cr)
	assert.Equal(s.T(), id, cr.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryCreate)).
		WithArgs(cr.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(categoryId))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreate)).
		WithArgs(cr.ID, categoryId, cr.Description, app.DefaultImage, app.DefaultImage).WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryCreate)).
		WithArgs(cr.Category).
		WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategoryCreate)).
		WithArgs(cr.Category).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), IncorrectCategory, err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreator() {
	cr := models.TestCreatorWithAwards()
	cr.ID = 1
	userId := int64(1)
	expected := *cr

	var awardsId sql.NullInt64
	awardsId.Valid = true
	awardsId.Int64 = cr.AwardsId
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname", "awards_id"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname, awardsId))
	get, err := s.repo.GetCreator(userId, expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)

	awardsId.Valid = false
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname", "awards_id"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname, awardsId))
	expected.AwardsId = rp.NoAwards
	get, err = s.repo.GetCreator(userId, expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).WillReturnError(sql.ErrNoRows)

	_, err = s.repo.GetCreator(userId, expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NotFound, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).WillReturnError(models.BDError)

	_, err = s.repo.GetCreator(userId, expected.ID)
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
	preapareCustomRows := CustomRows(*preapareRows)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(preapareCustomRows.Copy())
	get, err := s.repo.GetCreators()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), creators, get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(preapareCustomRows.Copy().RowError(0, models.BDError))
	_, err = s.repo.GetCreators()
	assert.EqualError(s.T(), err, repository.NewDBError(models.BDError).Error())

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(""))
	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).WillReturnError(models.BDError)
	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).WillReturnError(models.BDError)
	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateAvatar() {
	avatar := "d"
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	err := s.repo.UpdateAvatar(creatorId, avatar)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).WillReturnError(sql.ErrNoRows)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), app.UnknownError, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).WillReturnError(models.BDError)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateCover() {
	avatar := "d"
	creatorId := int64(1)
	runFunc := func(input ...interface{}) (res []interface{}) {
		id, _ := input[1].(int64)
		cover, _ := input[0].(string)
		err := s.repo.UpdateCover(id, cover)
		return []interface{}{err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{avatar, creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: queryUpdateCover,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId),
						RowError:   nil,
						RowClose:   nil,
					},
					Args:    []driver.Value{avatar, creatorId},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "NoRows",
			Args: []interface{}{avatar, creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   queryUpdateCover,
					Err:     sql.ErrNoRows,
					Rows:    nil,
					Args:    []driver.Value{avatar, creatorId},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BDError",
			Args: []interface{}{avatar, creatorId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(models.BDError),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   queryUpdateCover,
					Err:     models.BDError,
					Rows:    nil,
					Args:    []driver.Value{avatar, creatorId},
					RunType: models.Query,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func testQueryAddPagination(tableName string, res int64) models.TestQuery {
	return models.TestQuery{
		Query: "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1",
		Err:   nil,
		Rows: &models.TestRow{
			ReturnRows: sqlmock.NewRows([]string{"n_live_tup"}).AddRow(res),
			RowError:   nil,
			RowClose:   nil,
		},
		Args:    []driver.Value{tableName},
		RunType: models.Query,
	}
}

func testQueryAddPaginationWithError(tableName string, err error) models.TestQuery {
	return models.TestQuery{
		Query: "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1",
		Err:   err,
		Rows: &models.TestRow{
			ReturnRows: nil,
			RowError:   nil,
			RowClose:   nil,
		},
		Args:    []driver.Value{tableName},
		RunType: models.Query,
	}
}

func (s *SuiteCreatorRepository) TestCreatorRepository_SearchCreators() {
	creators := models.TestCreators()
	searchString := "dor"
	preapareRows := sqlmock.NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"})
	for index, cr := range creators {
		cr.ID = int64(index)
		creators[index] = cr
		preapareRows.AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname)
	}
	preapareCustomRows := CustomRows(*preapareRows)
	pag := &models.Pagination{Limit: 10, Offset: 0}
	categories := []string{"другое", "пол"}
	queryWithCategory, args, err := sqlx.In(querySearchCreators+queryCategorySearchCreators,
		utils.StringsToLowerCase(categories))
	require.NoError(s.T(), err)
	queryWithCategory = postgresql_utilits.CustomRebind(4, queryWithCategory)

	runFunc := func(input ...interface{}) (res []interface{}) {
		pagin, _ := input[0].(*models.Pagination)
		search, _ := input[1].(string)
		if len(input) == 2 {
			first, second := s.repo.SearchCreators(pagin, search)
			return []interface{}{first, second}
		}
		cats, _ := input[2].([]string)
		first, second := s.repo.SearchCreators(pagin, search, cats...)
		return []interface{}{first, second}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{pag, searchString},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{creators},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPagination("search_creators", 5000),
				{
					Query: querySearchCreators,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: preapareCustomRows.Copy(),
						RowError:   nil,
						RowClose:   nil,
					},
					Args:    []driver.Value{searchString, pag.Limit, pag.Offset},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "CorrectWithCategory",
			Args: []interface{}{pag, searchString, categories},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{creators},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPagination("search_creators", 5000),
				{
					Query: queryWithCategory,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: preapareCustomRows.Copy(),
						RowError:   nil,
						RowClose:   nil,
					},
					Args:    append([]driver.Value{searchString, pag.Limit, pag.Offset}, args[0], args[1]),
					RunType: models.Query,
				},
			},
		},
		{
			Name: "RowError",
			Args: []interface{}{pag, searchString},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]models.Creator(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPagination("search_creators", 5000),
				{
					Query: querySearchCreators,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: preapareCustomRows.Copy(),
						RowError: &models.TestRowError{
							Row: 2,
							Err: repository.DefaultErrDB,
						},
						RowClose: nil,
					},
					Args:    []driver.Value{searchString, pag.Limit, pag.Offset},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "ScanError",
			Args: []interface{}{pag, searchString},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{[]models.Creator(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPagination("search_creators", 5000),
				{
					Query: querySearchCreators,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"count"}).AddRow(1),
						RowError:   nil,
						RowClose:   nil,
					},
					Args:    []driver.Value{searchString, pag.Limit, pag.Offset},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BdErrorInBaseQuery",
			Args: []interface{}{pag, searchString},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]models.Creator(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPagination("search_creators", 5000),
				{
					Query:   querySearchCreators,
					Err:     repository.DefaultErrDB,
					Rows:    nil,
					Args:    []driver.Value{searchString, pag.Limit, pag.Offset},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BdErrorInPagQuery",
			Args: []interface{}{pag, searchString},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]models.Creator(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				testQueryAddPaginationWithError("search_creators", repository.DefaultErrDB),
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteCreatorRepository) TestCreatorRepository_ExistCreator() {
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	check, err := s.repo.ExistsCreator(creatorId)
	assert.NoError(s.T(), err)
	assert.True(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).WillReturnError(sql.ErrNoRows)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), app.UnknownError, err)
	assert.False(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).WillReturnError(models.BDError)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
	assert.False(s.T(), check)
}

func TestCreatorRepository(t *testing.T) {
	suite.Run(t, new(SuiteCreatorRepository))
}
