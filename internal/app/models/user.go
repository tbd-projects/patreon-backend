package models

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_LOGIN_LENGTH    = 5
	MAX_LOGIN_LENGTH    = 25
	MIN_NICKNAME_LENGTH = 5
	MAX_NICKNAME_LENGTH = 25
	MIN_PASSWORD_LENGTH = 6
	MAX_PASSWORD_LENGTH = 50
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
	HaveCreator       bool   `json:"have_creator"`
}

func (u *User) String() string {
	return fmt.Sprintf("{ID: %s, Login: %s}", strconv.Itoa(int(u.ID)), u.Login)
}

// Validate Errors:
//		IncorrectEmailOrPassword
//		IncorrectNickname
// Important can return some other error
func (u *User) Validate() error {
	err := validation.Errors{
		"login": validation.Validate(u.Login, validation.Required, validation.Length(MIN_LOGIN_LENGTH, MAX_LOGIN_LENGTH)),
		"password": validation.Validate(u.Password, validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(MIN_PASSWORD_LENGTH, MAX_PASSWORD_LENGTH)),
		"nickname": validation.Validate(u.Nickname, validation.Required, validation.Length(MIN_NICKNAME_LENGTH, MAX_NICKNAME_LENGTH)),
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
