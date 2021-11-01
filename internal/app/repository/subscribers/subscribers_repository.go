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
func (repo *SubscribersRepository) Create(subscriber *models.Subscriber, awardName string) error {
	queryAwardPrice := "SELECT price FROM awards WHERE creator_id = $1 AND name = $2"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id) VALUES ($1, $2)"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id) VALUES($1, $2, $3)"

	price := 0
	if err := repo.store.QueryRow(queryAwardPrice, subscriber.CreatorID, awardName).Scan(&price); err != nil {
		return repository.NewDBError(err)
	}
	begin, err := repo.store.Begin()
	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err := repo.store.QueryRow(queryAddPayment, price,
		subscriber.CreatorID, subscriber.UserID); err.Err() != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err.Err())
	}
	if err := repo.store.QueryRow(queryAddSubscribe,
		subscriber.UserID, subscriber.CreatorID); err.Err() != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err.Err())
	}

	if err = begin.Commit(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
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

// Get Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) Get(userID int64, creatorID int64) (bool, error) {
	query := "SELECT count(*) from subscribers where users_id = $1 and creator_id = $2"
	cnt := 0
	if res := repo.store.QueryRow(query, userID, creatorID).Scan(&cnt); res != nil {
		return false, repository.NewDBError(res)
	}
	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// Delete Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) Delete(subscriber *models.Subscriber) error {
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2"
	row, err := repo.store.Query(query, subscriber.UserID, subscriber.CreatorID)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}
