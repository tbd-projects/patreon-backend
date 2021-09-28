package models

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type User struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	Nickname          string `json:"nickname"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:",omitempty"`
	Avatar            string `json:"avatar,omitempty"`
}

func (u *User) String() string {
	return fmt.Sprintf("{ID: %s, Login: %s}", strconv.Itoa(u.ID), u.Login)
}
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Login, validation.Required, validation.Length(5, 25)),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 50)),
	)
}
func (u *User) MakePrivateDate() {
	u.Password = ""
	u.EncryptedPassword = ""
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := u.encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func (u *User) encryptString(s string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(enc), nil
}
