package repository_postgresql

import (
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/repository"
)

const (
	countCreatorPosts       = "SELECT count(*) as cnt from posts where creator_id = $1;"
	countCreatorSubscribers = "SELECT count(*) as cnt from subscribers where creator_id = $1;"
	//countCreatorPostsViews     = "SELECT sum(views) FROM posts GROUP BY creator_id HAVING creator_id = $1;"
	countCreatorPostsLastViews = "SELECT SUM(views) FROM (" +
		"SELECT posts.views, creator_id FROM POSTS" +
		"WHERE date_part('day', current_date::timestamptz - posts.date) < $2" +
		"GROUP BY creator_i HAVING creator_id = $1;"
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
func (r *StatisticsRepository) GetCountCreatorViews(creatorID int64, days int64) (int64, error) {
	var cnt int64
	err := r.store.QueryRow(countCreatorPostsLastViews, creatorID).Scan(&cnt)

	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return cnt, nil

}
