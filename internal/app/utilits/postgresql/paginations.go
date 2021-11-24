package postgresql_utilits

import (
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

const (
	queryStat = "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
)

// AddPagination Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func AddPagination(tableName string, pag *models.Pagination, db *sqlx.DB) (limit int64, offset int64, err error) {
	var numberRows int64
	if err = db.QueryRow(queryStat, tableName).Scan(&numberRows); err != nil {
		return app.InvalidInt, app.InvalidInt, repository.NewDBError(err)
	}

	numberRows -= pag.Limit
	if pag.Offset < numberRows {
		numberRows = pag.Offset
	}
	if numberRows < 0 {
		numberRows = 0
	}
	return pag.Limit, numberRows, nil
}
