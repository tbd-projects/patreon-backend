package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"patreon/internal/app"
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
func (u *User) Validate() error {
	validRes := validation.Errors{
		"login": validation.Validate(u.Login, validation.Required, validation.Length(5, 25)),
		"password": validation.Validate(u.Password, validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(6, 50)),
	}.Filter()

	var mapOfErr map[string]string
	er, ok := json.Marshal(validRes)
	if ok != nil {
		return InternalError
	}
	ok = json.Unmarshal(er, &mapOfErr)
	if ok != nil {
		return InternalError
	}

	if userValidError()("login") != nil || userValidError()("password") != nil {
		res := app.GeneralError{
			Err: IncorrectEmailOrPassword,
		}
		if err, ok := mapOfErr["login"]; ok {
			res.ExternalErr = errors.New(err)
		} else if err, ok := mapOfErr["password"]; ok {
			res.ExternalErr = errors.New(err)
		}
		return res
	}
	return nil
}
func (u *User) MakePrivateDate() {
	u.Password = ""
	u.EncryptedPassword = ""
}

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
