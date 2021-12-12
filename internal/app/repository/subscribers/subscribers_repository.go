package repository_subscribers

import (
	"github.com/jmoiron/sqlx"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type SubscribersRepository struct {
	store *sqlx.DB
}

func NewSubscribersRepository(store *sqlx.DB) *SubscribersRepository {
	return &SubscribersRepository{
		store: store,
	}
}

// Create Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) Create(subscriber *models.Subscriber, payToken string) error {
	queryAwardPrice := "SELECT price FROM awards WHERE awards_id = $1"
	queryAddPayment := "INSERT INTO payments(amount, creator_id, users_id, pay_token) VALUES($1, $2, $3, $4)"
	queryAddSubscribe := "INSERT INTO subscribers(users_id, creator_id, awards_id) VALUES ($1, $2, $3)"

	price := 0
	if err := repo.store.QueryRow(queryAwardPrice, subscriber.AwardID).Scan(&price); err != nil {
		return repository.NewDBError(err)
	}

	begin, err := repo.store.Begin()
	if err != nil {
		return repository.NewDBError(err)
	}

	row, err := begin.Query(queryAddPayment, price,
		subscriber.CreatorID, subscriber.UserID, payToken)

	if err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if row, err = begin.Query(queryAddSubscribe,
		subscriber.UserID, subscriber.CreatorID, subscriber.AwardID); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		_ = begin.Rollback()
		return repository.NewDBError(err)
	}

	if err = begin.Commit(); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// GetCreators Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) GetCreators(userID int64) ([]models.CreatorSubscribe, error) {
	queryCount := "SELECT count(*) as cnt from subscribers WHERE users_id = $1"
	querySelect := `
	SELECT DISTINCT s.creator_id, s.awards_id, category, description, nickname, cp.avatar, cover
	FROM subscribers s JOIN creator_profile cp ON s.creator_id = cp.creator_id
	JOIN users u ON cp.creator_id = u.users_id where s.users_id = $1
	`

	count := 0
	if err := repo.store.QueryRow(queryCount, userID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}

	var res []models.CreatorSubscribe

	rows, err := repo.store.Query(querySelect, userID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}
	var cur models.CreatorSubscribe
	for rows.Next() {
		if err = rows.Scan(&cur.ID, &cur.AwardsId, &cur.Category, &cur.Description, &cur.Nickname,
			&cur.Avatar, &cur.Cover); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		res = append(res, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil

}

// GetSubscribers Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) GetSubscribers(creatorID int64) ([]models.User, error) {
	queryCount := "SELECT count(*) as cnt from subscribers WHERE creator_id = $1"
	querySelect := `
	SELECT DISTINCT s.users_id, nickname, avatar
	from subscribers s join users u on s.users_id = u.users_id WHERE s.creator_id = $1`

	count := 0
	if err := repo.store.QueryRow(queryCount, creatorID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}

	res := make([]models.User, 0, count)

	rows, err := repo.store.Query(querySelect, creatorID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}
	var cur models.User
	for rows.Next() {
		if err = rows.Scan(&cur.ID, &cur.Nickname, &cur.Avatar); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		res = append(res, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// Get Errors:
//		app.GeneralError with Errors
//			repository.DefaultErrDB
func (repo *SubscribersRepository) Get(subscriber *models.Subscriber) (bool, error) {
	query := "SELECT count(*) as cnt from subscribers where users_id = $1 and creator_id = $2"
	cnt := 0
	if res := repo.store.QueryRow(query, subscriber.UserID, subscriber.CreatorID).Scan(&cnt); res != nil {
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
	query := "DELETE from subscribers where users_id = $1 and creator_id = $2 and awards_id = $3"
	row, err := repo.store.Query(query, subscriber.UserID, subscriber.CreatorID, subscriber.AwardID)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}
