package store

import (
	"os"
	"testing"
)

var dbUrl string

func TestMain(m *testing.M) {
	dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "host=localhost dbname=restapi_test sslmode=disable"
	}
	os.Exit(m.Run())
}
