package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/lib/pq"
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
	create_users_tb := "CREATE TABLE IF NOT EXISTS users" +
		"(" +
		"user_id            bigserial not null primary key," +
		"login              text      not null unique," +
		"nickname           text      not null unique," +
		"encrypted_password text      not null," +
		"avatar             text      not null" +
		");"
	create_creator_profile_db := "CREATE TABLE IF NOT EXISTS creator_profile" +
		"(" +
		"creator_id  integer not null primary key," +
		"category    text      not null," +
		"description text      not null," +
		"avatar      text      not null," +
		"cover       text      not null," +
		"foreign key (creator_id) references users(user_id) on delete cascade" +
		");"

	db.QueryRow(create_users_tb)
	db.QueryRow(create_creator_profile_db)

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
