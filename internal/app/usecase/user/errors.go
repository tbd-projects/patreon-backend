package usercase_user

import (
	"errors"
)

var (
	UserExist                = errors.New("user already exist")
	BadEncrypt               = errors.New("unsuccessful encrypt user")
	IncorrectEmailOrPassword = errors.New("incorrect email or password")
)
