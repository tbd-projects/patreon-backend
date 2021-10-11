package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	user := TestUser()
	assert.NoError(t, user.Encrypt())
	assert.NotEmpty(t, user.EncryptedPassword)
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		u       func() *User
		isValid bool
	}{
		{
			name: "valid user",
			u: func() *User {
				return TestUser()
			},
			isValid: true,
		},
		{
			name: "invalid login",
			u: func() *User {
				u := TestUser()
				u.Login = "A"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *User {
				u := TestUser()
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *User {
				u := TestUser()
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "with encrypted password and empty password",
			u: func() *User {
				u := TestUser()
				u.Password = ""
				u.EncryptedPassword = "encrypted"
				return u
			},
			isValid: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.isValid {
				assert.NoError(t, test.u().Validate())
			} else {
				assert.Error(t, test.u().Validate())
			}
		})
	}
}

func TestUser_MakePrivateDate(t *testing.T) {
	tests := []struct {
		name string
		u    *User
	}{
		{
			name: "Valid test",
			u:    TestUser(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.u.MakePrivateDate()
			assert.Empty(t, test.u.Password)
			assert.Empty(t, test.u.EncryptedPassword)
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		password string
		isValid  bool
	}{
		{
			name:     "invalidPassword",
			user:     User{},
			password: "pswd",
			isValid:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.user.Password = test.password
			assert.NoError(t, test.user.Encrypt())
			res := test.user.ComparePassword(test.password)
			if test.isValid {
				assert.True(t, res)
			} else {
				assert.False(t, res)
			}
		})
	}
}
