package sqlstore

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

var BDError = errors.New("BD error")

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
