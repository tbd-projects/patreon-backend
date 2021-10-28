package postgresql_utilits

import (
	"database/sql"
	"fmt"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type QueryInfo struct {
	Query string
	TableName string
}

func AddPagination(query *QueryInfo, pag *models.Pagination, db *sql.DB) (string, error) {
	queryStat := `SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1`

	var numberRows int64
	if err :=db.QueryRow(queryStat, query.TableName).Scan(&numberRows); err != nil {
		return "", repository.NewDBError(err)
	}

	numberRows -= pag.Limit
	if pag.Offset < numberRows {
		numberRows = pag.Offset
	}
	return fmt.Sprintf("%s LIMIT %d OFFEST %d", query.Query, pag.Limit, numberRows), nil
}
