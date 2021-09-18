package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (st *Store) Open() error {
	db, err := sql.Open("postgres", st.config.DataBaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	st.db = db

	return nil

}

func (st *Store) Close() error {
	err := st.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (st *Store) User() *UserRepository {
	if st.userRepository != nil {
		return st.userRepository
	}
	st.userRepository = NewUserRepository(st)

	return st.userRepository
}
