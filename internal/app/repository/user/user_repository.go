package repository_user

import (
	"database/sql"
	"patreon/internal/app/repository"
	"patreon/internal/models"

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

func (repo *UserRepository) Create(u *models.User) error {
	if err := repo.store.QueryRow("INSERT INTO users (login, nickname, encrypted_password, avatar) VALUES ($1, $2, $3, $4)"+
		"RETURNING user_id", u.Login, u.Nickname, u.EncryptedPassword, u.Avatar).Scan(&u.ID); err != nil {
		return parseDBError(err.(*pq.Error))
	}
	return nil
}

func (repo *UserRepository) FindByLogin(login string) (*models.User, error) {
	user := models.User{}

	if err := repo.store.QueryRow("SELECT user_id, login, encrypted_password from users where login=$1", login).
		Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, parseDBError(err.(*pq.Error))

	}

	return &user, nil
}
func (repo *UserRepository) FindByID(id int64) (*models.User, error) {
	user := models.User{}

	if err := repo.store.QueryRow("SELECT user_id, nickname, avatar from users where user_id=$1", id).
		Scan(&user.ID, &user.Nickname, &user.Avatar); err != nil {
		return nil, repository.NotFound
	}

	return &user, nil
}
