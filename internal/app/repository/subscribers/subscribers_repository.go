package repository_subscribers

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type SubscribersRepository struct {
	store *sql.DB
}

func NewSubscribersRepository(store *sql.DB) *SubscribersRepository {
	return &SubscribersRepository{
		store: store,
	}
}

// Create Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) Create(subscriber *models.Subscriber) error {
	query := "INSERT into subscribers(users_id, creator_id) VALUES ($1, $2)"
	if err := repo.store.QueryRow(query, subscriber.UserID, &subscriber.CreatorID); err.Err() != nil {
		return repository.NewDBError(err.Err())
	}
	return nil
}

// GetCreators Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) GetCreators(userID int64) ([]int64, error) {
	count := 0
	query := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	if err := repo.store.QueryRow(query, userID).Scan(&count); err != nil {
		return []int64{}, repository.NewDBError(err)
	}

	res := make([]int64, 0, count)

	query = "SELECT creator_id from subscribers WHERE users_id = $1"
	rows, err := repo.store.Query(query, userID)
	if err != nil {
		return []int64{}, repository.NewDBError(err)
	}
	var cur models.Subscriber
	for rows.Next() {
		if err = rows.Scan(&cur.CreatorID); err != nil {
			return []int64{}, repository.NewDBError(err)
		}
		if err = rows.Err(); err != nil {
			return []int64{}, repository.NewDBError(err)
		}
		res = append(res, cur.CreatorID)
	}
	if err = rows.Close(); err != nil {
		return []int64{}, repository.NewDBError(err)
	}
	return res, nil

}

// GetSubscribers Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) GetSubscribers(creatorID int64) ([]int64, error) {
	count := 0
	query := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	if err := repo.store.QueryRow(query, creatorID).Scan(&count); err != nil {
		return []int64{}, repository.NewDBError(err)
	}

	res := make([]int64, 0, count)

	query = "SELECT users_id from subscribers WHERE creator_id = $1"
	rows, err := repo.store.Query(query, creatorID)
	if err != nil {
		return []int64{}, repository.NewDBError(err)
	}
	var cur models.Subscriber
	for rows.Next() {
		if err = rows.Scan(&cur.UserID); err != nil {
			return nil, repository.NewDBError(err)
		}
		if err = rows.Err(); err != nil {
			return nil, repository.NewDBError(err)
		}
		res = append(res, cur.UserID)
	}
	if err = rows.Close(); err != nil {
		return nil, repository.NewDBError(err)
	}
	return res, nil
}
