package repository_postgresql

import (
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/repository"
)

const (
	countCreatorPosts          = "SELECT count(*) as cnt from posts where creator_id = $1;"
	countCreatorSubscribers    = "SELECT count(*) as cnt from subscribers where creator_id = $1;"
	countCreatorPostsLastViews = "select coalesce((select sum(views) " +
		"from (select posts.views, creator_id from posts " +
		"where date_part('day', current_date::timestamptz - posts.date) < $2) as sum_viewes " +
		"group by sum_viewes.creator_id " +
		"having sum_viewes.creator_id = $1), 0);"
	totalCreatorIncomes = "select coalesce(" +
		"(select sum(amount) from (select amount, creator_id from payments " +
		"where date_part('day', current_date::timestamptz - payments.date) < $2) as last_payments " +
		"group by last_payments.creator_id having last_payments.creator_id = $1), 0) as sum_payments;"
)

type StatisticsRepository struct {
	store *sqlx.DB
}

func NewStatisticsRepository(st *sqlx.DB) *StatisticsRepository {
	return &StatisticsRepository{
		store: st,
	}
}

// GetCountCreatorPosts Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (r *StatisticsRepository) GetCountCreatorPosts(creatorID int64) (int64, error) {
	var cnt int64
	err := r.store.QueryRow(countCreatorPosts, creatorID).Scan(&cnt)

	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return cnt, nil
}

// GetCountCreatorSubscribers Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (r *StatisticsRepository) GetCountCreatorSubscribers(creatorID int64) (int64, error) {
	var cnt int64
	err := r.store.QueryRow(countCreatorSubscribers, creatorID).Scan(&cnt)

	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return cnt, nil
}

// GetCountCreatorViews Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (r *StatisticsRepository) GetCountCreatorViews(creatorID int64, days int64) (int64, error) {
	var cnt int64
	err := r.store.QueryRow(countCreatorPostsLastViews, creatorID, days).Scan(&cnt)

	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return cnt, nil
}

// GetTotalIncome Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (r *StatisticsRepository) GetTotalIncome(creatorID int64, days int64) (float64, error) {
	var sum float64

	err := r.store.QueryRow(totalCreatorIncomes, creatorID, days).Scan(&sum)

	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return sum, nil
}
