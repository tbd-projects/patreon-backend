package sqlstore

import (
	"database/sql"
	"patreon/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	creatorRepository *CreatorRepository
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
func (st *Store) Creator() store.CreatorRepository {
	if st.creatorRepository != nil {
		return st.creatorRepository
	}
	st.creatorRepository = NewCreatorRepository(st)

	return st.creatorRepository
}
