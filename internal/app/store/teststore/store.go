package teststore

import (
	"patreon/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

func (st *Store) User() store.UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}
	st.userRepository = NewUserRepository(st)

	return st.userRepository
}
