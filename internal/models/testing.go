package models

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Login:    "student1999",
		Password: "1!2!3!",
	}

}
