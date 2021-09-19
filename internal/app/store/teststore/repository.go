package teststore

import (
	"patreon/internal/app/store"
	"patreon/internal/models"
)

type UserRepository struct {
	store *Store
	users map[string]*models.User
}

func NewUserRepository(st *Store) *UserRepository {
	return &UserRepository{
		store: st,
		users: map[string]*models.User{},
	}
}
func (repo *UserRepository) Create(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	repo.users[u.Login] = u
	u.ID = len(repo.users)

	return nil
}

func (repo *UserRepository) FindByLogin(login string) (*models.User, error) {
	u, ok := repo.users[login]
	if !ok {
		return nil, store.NotFound
	}

	return u, nil

}
