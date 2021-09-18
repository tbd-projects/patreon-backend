package store

import (
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, dbUrl string) (*Store, func(tables ...string)) {
	t.Helper()

	config := NewConfig()
	config.DataBaseUrl = dbUrl

	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	create_users_rb := "CREATE TABLE IF NOT EXISTS users (" +
		"user_id bigserial not null primary key," +
		"login text not null unique," +
		"password text not null," +
		"avatar text not null);"

	s.db.QueryRow(create_users_rb)

	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE",
				strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}
	}
}
