package repository_postgresql

import (
	"fmt"
	"patreon/internal/app/models"
	db_models "patreon/internal/app/models"
	"patreon/internal/app/repository"
	putilits "patreon/internal/app/utilits/postgresql"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
)

const (
	querySelectUserPayments = "SELECT p.amount, p.date, p.creator_id, u.nickname, cp.category, cp.description FROM payments p " +
		"JOIN creator_profile cp on p.creator_id = cp.creator_id " +
		"JOIN users u on cp.creator_id = u.users_id where p.users_id = $1 " +
		"ORDER BY p.date DESC "

	querySelectCreatorPayments = "SELECT p.amount, p.date, p.users_id, u.nickname FROM payments p " +
		"JOIN users u on p.users_id = u.users_id where p.creator_id = $1 " +
		"ORDER BY p.date DESC "
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
func (repo *PaymentsRepository) GetUserPayments(userID int64, pag *db_models.Pagination) ([]models.UserPayments, error) {
	query := querySelectUserPayments

	limit, offset, err := putilits.AddPagination("payments", pag, repo.store)
	query = query + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	if limit == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(query, userID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.UserPayments, 0, limit)

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

// GetCreatorPayments Errors:
//		repository.NotFound
//		app.GeneralError with Errors:
//			repository.DefaultErrDB
func (repo *PaymentsRepository) GetCreatorPayments(creatorID int64, pag *db_models.Pagination) ([]models.CreatorPayments, error) {
	query := querySelectCreatorPayments

	limit, offset, err := putilits.AddPagination("payments", pag, repo.store)
	query = query + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	if err != nil {
		return nil, err
	}
	if limit == 0 {
		return nil, repository.NotFound
	}

	rows, err := repo.store.Query(query, creatorID)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	paymentsRes := make([]models.CreatorPayments, 0, limit)

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
