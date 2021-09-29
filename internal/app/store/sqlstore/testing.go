package sqlstore

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *Suite) InitBD() {
	s.T().Helper()

	var err error
	s.DB, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
}
