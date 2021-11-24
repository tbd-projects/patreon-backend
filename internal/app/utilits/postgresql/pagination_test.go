package postgresql_utilits

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhashkevych/go-sqlxmock"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func MockCallAddPagination(mock sqlmock.Sqlmock, tableName string, res int64) {
	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(res))
}

func MockCallAddPaginationWithError(mock sqlmock.Sqlmock, tableName string, err error) {
	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnError(err)
}

func Test_AddPagintation(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}

	tableName := "dddd"
	pag := &models.Pagination{Limit: 10, Offset: 20}
	MockCallAddPagination(mock, tableName, 5000)
	limit, offset, err := AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, pag.Offset)

	MockCallAddPagination(mock, tableName, 10)
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, int64(0))

	MockCallAddPagination(mock, tableName, 5)
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, int64(0))

	MockCallAddPaginationWithError(mock, tableName, repository.DefaultErrDB)
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.Error(t, err, repository.NewDBError(repository.DefaultErrDB))

	require.NoError(t, mock.ExpectationsWereMet())
}
