package postgresql

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

	"github.com/pkg/errors"
)

type PaymentsRepository struct {
	store *sql.DB
}

func NewPaymentsRepository(store *sql.DB) *PaymentsRepository {
	return &PaymentsRepository{
		store: store,
	}
}

// GetUserPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetUserPayments(userID int64) ([]models.Payment, error) {
	queryCount := "SELECT count(*) FROM payments where users_id = $1"
	querySelect := "SELECT amount, date, creator_id FROM payments where users_id = $1"

	count := 0

	if err := repo.store.QueryRow(queryCount, userID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	if count == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(querySelect)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.Payment, 0, count)

	for rows.Next() {
		cur := models.Payment{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.CreatorID); err != nil {
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

// GetCreatorPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetCreatorPayments(creatorID int64) ([]models.Payment, error) {
	queryCount := "SELECT count(*) FROM payments where creator_id = $1"
	querySelect := "SELECT amount, date, users_id FROM payments where creator_id = $1"

	count := 0

	if err := repo.store.QueryRow(queryCount, creatorID).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	if count == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(querySelect)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.Payment, 0, count)

	for rows.Next() {
		cur := models.Payment{}
		if err = rows.Scan(&cur.Amount, &cur.Date, &cur.UserID); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(errors.Wrapf(err, "method - GetCreatorPayments"+
				"invalid data in db: table payments"))
		}
		paymentsRes = append(paymentsRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return paymentsRes, nil
}
