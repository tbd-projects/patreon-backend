package store

import (
	"os"
	"testing"
)

var dbUrl string

func TestMain(m *testing.M) {
	dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "host=localhost dbname=restapi_test sslmode=disable user=postgres password=1111"
	}
	os.Exit(m.Run())
}
