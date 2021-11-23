package postgresql_utilits

import (
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func Test_AddPagintation(t *testing.T) {
	queryStat := "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}

	tableName := "dddd"
	pag := &models.Pagination{Limit: 10, Offset: 20}
	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	limit, offset, err := AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, pag.Offset)

	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(10)))
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, int64(0))

	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5)))
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.NoError(t, err)
	assert.Equal(t, limit, pag.Limit)
	assert.Equal(t, offset, int64(0))

	mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnError(repository.DefaultErrDB)
	limit, offset, err = AddPagination(tableName, pag, db)
	assert.Error(t, err, repository.NewDBError(repository.DefaultErrDB))
}
