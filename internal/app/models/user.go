package models

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type User struct {
	ID                int64  `json:"id"`
	Login             string `json:"login"`
	Nickname          string `json:"nickname"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:",omitempty"`
	Avatar            string `json:"avatar,omitempty"`
}

func (u *User) String() string {
	return fmt.Sprintf("{ID: %s, Login: %s}", strconv.Itoa(int(u.ID)), u.Login)
}

// Validate Errors:
//		IncorrectEmailOrPassword
// Important can return some other error
func (u *User) Validate() error {
	err := validation.Errors{
		"login": validation.Validate(u.Login, validation.Required, validation.Length(5, 25)),
		"password": validation.Validate(u.Password, validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(6, 50)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(err, "failed error getting in validate user")
	}

	if knowError = extractValidateError(userValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

func (u *User) MakeEmptyPassword() {
	u.Password = ""
	u.EncryptedPassword = ""
}

// Encrypt Errors:
// 		EmptyPassword
// Important can return some other error
func (u *User) Encrypt() error {
	if len(u.Password) == 0 {
		return EmptyPassword
	}
	enc, err := u.encryptString(u.Password)
	if err != nil {
		return err
	}
	u.EncryptedPassword = enc
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
