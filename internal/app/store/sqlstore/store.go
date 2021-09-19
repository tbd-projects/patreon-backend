package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"patreon/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (st *Store) User() store.UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}
	st.userRepository = NewUserRepository(st)

	return st.userRepository
}
