package repository_postgresql

import (
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	queryCountUserPayments  = "SELECT count(*) FROM payments where users_id = $1"
	querySelectUserPayments = "SELECT p.amount, p.date, p.creator_id, u.nickname, cp.category, cp.description FROM payments p " +
		"JOIN creator_profile cp on p.creator_id = cp.creator_id " +
		"JOIN users u on cp.creator_id = u.users_id where p.users_id = $1"

	queryCountCreatorPayments  = "SELECT count(*) FROM payments where creator_id = $1"
	querySelectCreatorPayments = "SELECT p.amount, p.date, p.users_id, u.nickname FROM payments p " +
		"JOIN users u on p.users_id = u.users_id where p.creator_id = $1"
)

type PaymentsRepository struct {
	store *sqlx.DB
}

func NewPaymentsRepository(store *sqlx.DB) *PaymentsRepository {
	return &PaymentsRepository{
		store: store,
	}
}

// GetUserPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetUserPayments(userID int64) ([]models.UserPayments, error) {
	count := 0

	if err := repo.store.QueryRow(queryCountUserPayments, userID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	if count == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(querySelectUserPayments, userID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.UserPayments, 0, count)

	for rows.Next() {
		cur := models.UserPayments{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.CreatorID,
			&cur.CreatorNickname, &cur.CreatorCategory, &cur.CreatorDescription); err != nil {

			_ = rows.Close()
			return nil, repository.NewDBError(errors.Wrapf(err, "method - GetUserPayments"+
				"invalid data in db: table payments"))
		}
		paymentsRes = append(paymentsRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return paymentsRes, nil
}

func (repo *PaymentsRepository) GetCreatorPayments(creatorID int64) ([]models.CreatorPayments, error) {

	count := 0

	if err := repo.store.QueryRow(queryCountCreatorPayments, creatorID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	if count == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(querySelectCreatorPayments, creatorID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.CreatorPayments, 0, count)

	for rows.Next() {
		cur := models.CreatorPayments{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.UserID, &cur.UserNickname); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(errors.Wrapf(err, "method - GetUserPayments"+
				"invalid data in db: table payments"))
		}
		paymentsRes = append(paymentsRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return paymentsRes, nil
}
