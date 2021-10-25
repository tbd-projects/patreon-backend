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

	if err := repo.store.QueryRow("SELECT users_id, login, nickname, avatar, encrypted_password "+
		"from users where users_id=$1", id).
		Scan(&user.ID, &user.Login, &user.Nickname, &user.Avatar, &user.EncryptedPassword); err != nil {
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
	if err := repo.store.QueryRow("UPDATE users SET encrypted_password = $1"+
		"WHERE users_id = $2", newEncryptedPassword, id).Scan(); err != nil {
		return err
	}
	return nil
}

// UpdateAvatar Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *UserRepository) UpdateAvatar(id int64, newAvatar string) error {
	if err := repo.store.QueryRow("UPDATE users SET avatar = $1"+
		"WHERE users_id = $2", newAvatar, id).Scan(); err != nil {
		return err
	}
	return nil
}
