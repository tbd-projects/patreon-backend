package store

import (
	"fmt"
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

func (repo *UserRepository) Create(u *models.User) (*models.User, error) {
	if err := repo.store.db.QueryRow("INSERT INTO users (login, password, avatar) VALUES ($1, $2, $3)"+
		"RETURNING ID", u.Login, u.Password, u.Avatar).Scan(&u.ID); err != nil {
		return nil, err
	}
	return u, nil
}
func (repo *UserRepository) FindByLogin(login string) (*models.User, error) {
	user := models.User{}
	if err := repo.store.db.QueryRow(fmt.Sprintf("SELECT user_id, login, password from users where login = $1", login)).
		Scan(&user.ID, &user.Login, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}
