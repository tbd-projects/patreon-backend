package repository_user

import (
	"database/sql"
	"github.com/lib/pq"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
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
	if err := repo.store.QueryRow("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4)"+
		"RETURNING user_id", u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).Scan(&u.ID); err != nil {
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

	if err := repo.store.QueryRow("SELECT user_id, login, encrypted_password from users where login=$1", login).
		Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
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

	if err := repo.store.QueryRow("SELECT user_id, nickname, avatar from users where user_id=$1", id).
		Scan(&user.ID, &user.Nickname, &user.Avatar); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return &user, nil
}
