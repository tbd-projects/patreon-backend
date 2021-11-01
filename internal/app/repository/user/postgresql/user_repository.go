package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"

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
	if err := repo.store.QueryRow("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4) "+
		"RETURNING users_id", u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).Scan(&u.ID); err != nil {
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
	user := models.User{}

	if err := repo.store.QueryRow("SELECT users_id, login, nickname, avatar, encrypted_password "+
		"from users where login=$1", login).
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

// UpdatePassword Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdatePassword(id int64, newEncryptedPassword string) error {
	row, err := repo.store.Query("UPDATE users SET encrypted_password = $1 WHERE users_id = $2",
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
	row, err := repo.store.Query("UPDATE users SET avatar = $1 WHERE users_id = $2", newAvatar, id)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}
