package sqlstore

import (
	"patreon/internal/app/store"
	"patreon/internal/models"
)

type UserRepository struct {
	store *Store
}

func NewUserRepository(st *Store) *UserRepository {
	return &UserRepository{
		store: st,
	}
}

func (repo *UserRepository) Create(u *models.User) error {
	if err := repo.store.db.QueryRow("INSERT INTO users (login, encrypted_password, avatar) VALUES ($1, $2, $3)"+
		"RETURNING user_id", u.Login, u.EncryptedPassword, u.Avatar).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) FindByLogin(login string) (*models.User, error) {
	user := models.User{}

	if err := repo.store.db.QueryRow("SELECT user_id, login, encrypted_password from users where login=$1", login).
		Scan(&user.ID, &user.Login, &user.EncryptedPassword); err != nil {
		return nil, store.NotFound
	}

	return &user, nil
}
func (repo *UserRepository) FindByID(id int64) (*models.User, error) {
	user := models.User{}

	if err := repo.store.db.QueryRow("SELECT user_id, login, avatar from users where user_id=$1", id).
		Scan(&user.ID, &user.Login, &user.Avatar); err != nil {
		return nil, store.NotFound
	}

	return &user, nil
}
