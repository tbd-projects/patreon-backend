package models

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Login:    "student1999",
		Password: "1!2!3!",
	}

}
func TestCreator(t *testing.T) *Creator {
	t.Helper()

	return &Creator{
		ID:          1,
		Category:    "podcasts",
		Description: "i love podcasts",
	}

}
