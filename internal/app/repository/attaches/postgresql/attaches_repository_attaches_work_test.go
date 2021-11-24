package repository_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"time"
)

func getAttachTypeAndIdTestQuery() models.TestQuery {
	return models.TestQuery{
		Query: getDataTypeAndIdQuery,
		Err:   nil,
		Rows: &models.TestRow{
			ReturnRows: sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "2").
				AddRow(2, models.Text),
		},
		RunType: models.Query,
	}
}

func getAttachTypeAndIdTestQueryError(err error) models.TestQuery {
	return models.TestQuery{
		Query:   getDataTypeAndIdQuery,
		Err:     err,
		Rows:    nil,
		RunType: models.Query,
	}
}

func markUnusedAttachQuery(level int64, postId int64) models.TestQuery {
	return models.TestQuery{
		Query:   makeUnusedAttachQuery,
		Args:    []driver.Value{level, postId},
		Err:     nil,
		RunType: models.Exec,
	}
}

func markUnusedAttachQueryError(level int64, postId int64, err error) models.TestQuery {
	return models.TestQuery{
		Query:   makeUnusedAttachQuery,
		Args:    []driver.Value{level, postId},
		Err:     err,
		Rows:    nil,
		RunType: models.Exec,
	}
}

func deleteUnusedAttachQuery(level int64, postId int64) models.TestQuery {
	return models.TestQuery{
		Query:   deleteUnusedQuery,
		Args:    []driver.Value{level, postId},
		Err:     nil,
		RunType: models.Exec,
	}
}

func deleteUnusedAttachQueryError(level int64, postId int64, err error) models.TestQuery {
	return models.TestQuery{
		Query:   deleteUnusedQuery,
		Args:    []driver.Value{level, postId},
		Err:     err,
		Rows:    nil,
		RunType: models.Exec,
	}
}

func updateAttach(level int64, id int64) []models.TestQuery {
	return []models.TestQuery{
		getAttachTypeAndIdTestQuery(),
		{
			Query: updateAttachFilesQuery,
			Args:  []driver.Value{level, id},
			Err:   nil,
			Rows: &models.TestRow{
				ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(id),
			},
			RunType: models.Query,
		},
	}
}

func createAttachError(err error) models.TestQuery {
	return getAttachTypeAndIdTestQueryError(err)
}

func createAttach(postId int64, typeId int64, attaches models.Attach) []models.TestQuery {
	query := "INSERT INTO posts_data (post_id, type, data, level) VALUES (?, ?, ?, ?) RETURNING data_id"
	return []models.TestQuery{
		getAttachTypeAndIdTestQuery(),
		{
			Query: query,
			Err:   nil,
			Args:  []driver.Value{postId, typeId, attaches.Value, attaches.Level},
			Rows: &models.TestRow{
				ReturnRows: sqlmock.NewRows([]string{"id"}).
					AddRow(attaches.Id),
			},
			RunType: models.Query,
		},
	}
}

func updateAttachError(level int64, id int64, err error) models.TestQuery {
	return models.TestQuery{
		Query: updateAttachFilesQuery,
		Args:  []driver.Value{level, id},
		Err:   err,
		RunType: models.Query,
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_getAttachTypeAndId() {
	attachTypes := map[models.DataType]int64{"2": 1, "3": 2}
	runFunc := func(input ...interface{}) (res []interface{}) {
		s.repo.dataTypes = nil
		s.repo.lastUpdate = time.Unix(1, 1)
		values, err := s.repo.getAttachTypeAndId()
		return []interface{}{values, err}
	}

	runFuncHaveTypes := func(input ...interface{}) (res []interface{}) {
		s.repo.dataTypes = attachTypes
		s.repo.lastUpdate = time.Now()
		values, err := s.repo.getAttachTypeAndId()
		return []interface{}{values, err}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{attachTypes},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: getDataTypeAndIdQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "2").
							AddRow(2, "3"),
					},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "CorrectWithHaveDataTypes",
			Args: []interface{}{},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{attachTypes},
			},
			RunFunc: runFuncHaveTypes,
		},
		{
			Name: "RowError",
			Args: []interface{}{},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{map[models.DataType]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: getDataTypeAndIdQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id", "name"}).
							AddRow(1, "2").
							AddRow(2, "3"),
						RowError: &models.TestRowError{
							Row: 1,
							Err: repository.DefaultErrDB,
						},
					},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "ScanError",
			Args: []interface{}{},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedReturns: []interface{}{map[models.DataType]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query: getDataTypeAndIdQuery,
					Err:   nil,
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).
							AddRow("2"),
					},
					RunType: models.Query,
				},
			},
		},
		{
			Name: "BdError",
			Args: []interface{}{},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{map[models.DataType]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Query:   getDataTypeAndIdQuery,
					Err:     repository.DefaultErrDB,
					Rows:    nil,
					RunType: models.Query,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_markUnusedAttach() {
	postId := int64(5)
	runFunc := func(input ...interface{}) (res []interface{}) {
		id, _ := input[0].(int64)

		trans, err := s.repo.store.Beginx()
		require.NoError(s.T(), err)

		errRet := s.repo.markUnusedAttach(trans, id)

		err = trans.Commit()
		require.NoError(s.T(), err)

		return []interface{}{errRet}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{postId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query:   makeUnusedAttachQuery,
					Args:    []driver.Value{UnusedAttach, postId},
					Err:     nil,
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "BDError",
			Args: []interface{}{postId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query:   makeUnusedAttachQuery,
					Args:    []driver.Value{UnusedAttach, postId},
					Err:     repository.DefaultErrDB,
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_deleteUnused() {
	postId := int64(5)
	runFunc := func(input ...interface{}) (res []interface{}) {
		id, _ := input[0].(int64)

		trans, err := s.repo.store.Beginx()
		require.NoError(s.T(), err)

		errRet := s.repo.deleteUnused(trans, id)

		err = trans.Commit()
		require.NoError(s.T(), err)

		return []interface{}{errRet}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{postId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query:   deleteUnusedQuery,
					Args:    []driver.Value{UnusedAttach, postId},
					Err:     nil,
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "BDError",
			Args: []interface{}{postId},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				{
					Query:   deleteUnusedQuery,
					Args:    []driver.Value{UnusedAttach, postId},
					Err:     repository.DefaultErrDB,
					RunType: models.Exec,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_updateAttach() {
	textAttach := &models.Attach{
		Id:    1,
		Type:  models.Text,
		Value: "don",
		Level: 1,
	}

	fileAttach := &models.Attach{
		Id:    1,
		Type:  "2",
		Value: "don",
		Level: 1,
	}

	runFunc := func(input ...interface{}) (res []interface{}) {
		attach, _ := input[0].(*models.Attach)

		s.repo.dataTypes = nil
		s.repo.lastUpdate = time.Unix(1, 1)

		trans, err := s.repo.store.Beginx()
		require.NoError(s.T(), err)

		errRet := s.repo.updateAttach(trans, attach)

		err = trans.Commit()
		require.NoError(s.T(), err)

		return []interface{}{errRet}
	}

	testings := []models.TestCase{
		{
			Name: "CorrectText",
			Args: []interface{}{textAttach},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query: updateAttachQuery,
					Args:  []driver.Value{textAttach.Value, textAttach.Level, textAttach.Id},
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(fileAttach.Id),
					},
					Err:     nil,
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "CorrectFile",
			Args: []interface{}{fileAttach},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query: updateAttachFilesQuery,
					Args:  []driver.Value{fileAttach.Level, fileAttach.Id},
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).AddRow(fileAttach.Id),
					},
					Err:     nil,
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "NotFound",
			Args: []interface{}{fileAttach},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NotFound,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query:   updateAttachFilesQuery,
					Args:    []driver.Value{fileAttach.Level, fileAttach.Id},
					Err:     sql.ErrNoRows,
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "BDErrorQuery",
			Args: []interface{}{fileAttach},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query:   updateAttachFilesQuery,
					Args:    []driver.Value{fileAttach.Level, fileAttach.Id},
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "InvalidType",
			Args: []interface{}{&models.Attach{Type: "dut"}},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     UnknownDataFormat,
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "BDErrorGetTypes",
			Args: []interface{}{&models.Attach{Type: "dut"}},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQueryError(repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_ApplyChangeAttaches() {
	postId := int64(2)
	newAttaches := []models.Attach{
		{
			Id:    1,
			Type:  "2",
			Value: "2",
			Level: 2,
		},
	}
	updAttaches := []models.Attach{
		{
			Id:    3,
			Type:  "2",
			Value: "2",
			Level: 1,
		},
	}

	runFunc := func(input ...interface{}) (res []interface{}) {
		id, _ := input[0].(int64)
		newAtt, _ := input[1].([]models.Attach)
		updAtt, _ := input[2].([]models.Attach)
		s.repo.dataTypes = nil
		resIds, errRet := s.repo.ApplyChangeAttaches(id, newAtt, updAtt)

		return []interface{}{resIds, errRet}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{[]int64{3, 1}},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQuery(UnusedAttach, postId),
				createAttach(postId, 1, newAttaches[0])[0],
				createAttach(postId, 1, newAttaches[0])[1],
				updateAttach(updAttaches[0].Level, updAttaches[0].Id)[1],
				deleteUnusedAttachQuery(UnusedAttach, postId),
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "ErrorBegin",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     repository.DefaultErrDB,
					RunType: models.TransBegin,
				},
			},
		},
		{
			Name: "ErrorMarkUnused",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQueryError(UnusedAttach, postId, repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrorCreateAttaches",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQuery(UnusedAttach, postId),
				createAttachError(repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrorUpdateAttaches",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQuery(UnusedAttach, postId),
				createAttach(postId, 1, newAttaches[0])[0],
				createAttach(postId, 1, newAttaches[0])[1],
				updateAttachError(updAttaches[0].Level, updAttaches[0].Id, repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrorDeleteUnused",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQuery(UnusedAttach, postId),
				createAttach(postId, 1, newAttaches[0])[0],
				createAttach(postId, 1, newAttaches[0])[1],
				updateAttach(updAttaches[0].Level, updAttaches[0].Id)[1],
				deleteUnusedAttachQueryError(UnusedAttach, postId, repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransRollback,
				},
			},
		},
		{
			Name: "ErrorCommit",
			Args: []interface{}{postId, newAttaches, updAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]int64(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				markUnusedAttachQuery(UnusedAttach, postId),
				createAttach(postId, 1, newAttaches[0])[0],
				createAttach(postId, 1, newAttaches[0])[1],
				updateAttach(updAttaches[0].Level, updAttaches[0].Id)[1],
				deleteUnusedAttachQuery(UnusedAttach, postId),
				{
					Err:     repository.DefaultErrDB,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}

func (s *SuiteAttachesRepository) TestAttachesRepository_createAttaches() {
	query := "INSERT INTO posts_data (post_id, type, data, level) VALUES (?, ?, ?, ?) RETURNING data_id"
	newAttaches := []models.Attach{
		{
			Id:    20,
			Value: "dor",
			Level: 1,
			Type:  "2",
		},
	}
	postId := int64(5)
	runFunc := func(input ...interface{}) (res []interface{}) {
		id, _ := input[0].(int64)
		newAttach, _ := input[1].([]models.Attach)

		s.repo.dataTypes = nil
		s.repo.lastUpdate = time.Unix(1, 1)

		trans, err := s.repo.store.Beginx()
		require.NoError(s.T(), err)

		values, errRet := s.repo.createAttaches(trans, id, newAttach)

		err = trans.Commit()
		require.NoError(s.T(), err)

		return []interface{}{values, errRet}
	}

	testings := []models.TestCase{
		{
			Name: "Correct",
			Args: []interface{}{postId, newAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     nil,
				ExpectedReturns: []interface{}{newAttaches},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query: query,
					Err:   nil,
					Args:  []driver.Value{postId, 1, newAttaches[0].Value, newAttaches[0].Level},
					Rows: &models.TestRow{
						ReturnRows: sqlmock.NewRows([]string{"id"}).
							AddRow(20),
					},
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "DBErrorQuery",
			Args: []interface{}{postId, newAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]models.Attach(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Query:   query,
					Args:    []driver.Value{postId, 1, newAttaches[0].Value, newAttaches[0].Level},
					Err:     repository.DefaultErrDB,
					RunType: models.Query,
				},
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "DBErrorGetType",
			Args: []interface{}{postId, newAttaches},
			Expected: models.TestExpected{
				HaveError:       true,
				ExpectedErr:     repository.NewDBError(repository.DefaultErrDB),
				ExpectedReturns: []interface{}{[]models.Attach(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQueryError(repository.DefaultErrDB),
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
		{
			Name: "UnknownType",
			Args: []interface{}{postId,
				func() []models.Attach {
					dst := make([]models.Attach, len(newAttaches))
					copy(dst, newAttaches)
					dst[0].Type = "dor"
					return dst
				}(),
			},
			Expected: models.TestExpected{
				HaveError:       true,
				CheckError:      true,
				ExpectedReturns: []interface{}{[]models.Attach(nil)},
			},
			RunFunc: runFunc,
			Queries: []models.TestQuery{
				{
					Err:     nil,
					RunType: models.TransBegin,
				},
				getAttachTypeAndIdTestQuery(),
				{
					Err:     nil,
					RunType: models.TransCommit,
				},
			},
		},
	}

	for _, test := range testings {
		s.RunTestCase(test)
	}
}
