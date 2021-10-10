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
	//_, err := govalidator.ValidateStruct(*u)
	//if err != nil {
	//	errs := err.(govalidator.Errors).Errors()
	//	res := ""
	//	for _, e := range errs {
	//		res += e.Error() + " "
	//	}
	//	return errors.
	//}
	//return err
	//err := validation.Errors{
	//	"name":  validation.Validate(c.Name, validation.Required, validation.Length(5, 20)),
	//	"email": validation.Validate(c.Name, validation.Required, is.Email),
	//	"num":   validation.Validate(c.Num, validation.In(1, 2, 3)),
	return validation.ValidateStruct(u,
		validation.Field(&u.Login, validation.Required, validation.Length(5, 25)),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 50)),
	)
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
