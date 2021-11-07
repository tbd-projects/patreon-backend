package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

	"github.com/pkg/errors"

	"github.com/lib/pq"
)

type UserRepository struct {
	store *sql.DB
}

func NewUserRepository(st *sql.DB) *UserRepository {
	return &UserRepository{
		store: st,
	}
}

// Create Errors:
// 		LoginAlreadyExist
// 		NicknameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) RETURNING users_id`

	if err := repo.store.QueryRow(query, u.Login, u.Nickname, u.EncryptedPassword, app.DefaultImage).Scan(&u.ID); err != nil {
		if _, ok := err.(*pq.Error); ok {
			return parsePQError(err.(*pq.Error))
		}
		return repository.NewDBError(err)
	}
	return nil
}

// FindByLogin Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *UserRepository) FindByLogin(login string) (*models.User, error) {
	query := `SELECT users_id, login, nickname, avatar, encrypted_password from users where login=$1`
	user := models.User{}

	if err := repo.store.QueryRow(query, login).
		Scan(&user.ID, &user.Login, &user.Nickname, &user.Avatar, &user.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)

	}

	return &user, nil
}

// FindByID Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) FindByID(id int64) (*models.User, error) {
	user := models.User{}
	query := `SELECT users_id, login, nickname, users.avatar, encrypted_password, cp.creator_id IS NOT NULL
	from users LEFT JOIN creator_profile AS cp ON (users.users_id = cp.creator_id) where users_id=$1`

	if err := repo.store.QueryRow(query, id).
		Scan(&user.ID, &user.Login, &user.Nickname, &user.Avatar, &user.EncryptedPassword, &user.HaveCreator); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return &user, nil
}

// FindByNickname Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) FindByNickname(nickname string) (*models.User, error) {
	user := models.User{}
	query := `SELECT users_id, login, nickname, users.avatar, encrypted_password, cp.creator_id IS NOT NULL
	from users LEFT JOIN creator_profile AS cp ON (users.users_id = cp.creator_id) where nickname=$1`

	if err := repo.store.QueryRow(query, nickname).
		Scan(&user.ID, &user.Login, &user.Nickname, &user.Avatar, &user.EncryptedPassword, &user.HaveCreator); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return &user, nil
}

// UpdatePassword Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdatePassword(id int64, newEncryptedPassword string) error {
	query := `UPDATE users SET encrypted_password = $1 WHERE users_id = $2`

	row, err := repo.store.Query(query,
		newEncryptedPassword, id)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// UpdateAvatar Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdateAvatar(id int64, newAvatar string) error {
	query := `UPDATE users SET avatar = $1 WHERE users_id = $2`

	row, err := repo.store.Query(query, newAvatar, id)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// UpdateNickname Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdateNickname(id int64, newNickname string) error {
	query := `UPDATE users SET nickname = $1 WHERE users_id = $2`

	row, err := repo.store.Exec(query, newNickname, id)
	if err != nil {
		return repository.NewDBError(err)
	}
	if cntChangesRows, err := row.RowsAffected(); err != nil || cntChangesRows != 1 {
		if err != nil {
			return repository.NewDBError(err)
		}
		return repository.NewDBError(
			errors.Wrapf(err,
				"UPDATE_NICKNAME_REPO: expected changes only one row in db, change %v", cntChangesRows))
	}
	return nil
}
