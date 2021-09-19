package models_test

import (
	"github.com/stretchr/testify/assert"
	"patreon/internal/models"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	user := models.TestUser(t)
	assert.NoError(t, user.BeforeCreate())
	assert.NotEmpty(t, user.EncryptedPassword)
}
func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		u       func() *models.User
		isValid bool
	}{
		{
			name: "valid user",
			u: func() *models.User {
				return models.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "invalid login",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Login = "A"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *models.User {
				u := models.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "with encrypted password and empty password",
			u: func() *models.User {
				u := models.TestUser(t)
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
