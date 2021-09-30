package models

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Login:    "student1999",
		Password: "1!2!3!",
	}

}
func TestUsers(t *testing.T) []User {
	t.Helper()
	u1 := User{
		Login:    "test1",
		Nickname: "test1",
		Password: "123456",
	}
	u2 := User{
		Login:    "test2",
		Nickname: "test2",
		Password: "123456",
	}
	u3 := User{
		Login:    "test3",
		Nickname: "test3",
		Password: "123456",
	}
	return []User{u1, u2, u3}

}
func TestCreator(t *testing.T) *Creator {
	t.Helper()

	return &Creator{
		ID:          1,
		Category:    "podcasts",
		Description: "i love podcasts",
	}
}
func TestCreators(t *testing.T) []Creator {
	t.Helper()
	cr1 := Creator{
		ID:          1,
		Category:    "podcasts",
		Description: "i love podcasts",
	}
	cr2 := Creator{
		ID:          2,
		Category:    "blog",
		Description: "i love podcasts",
	}
	cr3 := Creator{
		ID:          3,
		Category:    "movies",
		Description: "i love podcasts",
	}
	return []Creator{cr1, cr2, cr3}

}
