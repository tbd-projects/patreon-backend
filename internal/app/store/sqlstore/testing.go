package sqlstore

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"testing"
)

func TestDB(t *testing.T, dbUrl string) (*sql.DB, func(tables ...string)) {
	t.Helper()

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	create_users_rb := "CREATE TABLE IF NOT EXISTS users (" +
		"user_id bigserial not null primary key," +
		"login text not null unique," +
		"encrypted_password text not null," +
		"avatar text not null);"

	db.QueryRow(create_users_rb)

	return db, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE",
				strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
