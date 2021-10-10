package usercase_user

import (
	"errors"
)

var (
	UserExist                = errors.New("user already exist")
	EmptyPassword            = errors.New("user password is empty")
	BadEncrypt               = errors.New("unsuccessful encrypt user")
	IncorrectEmailOrPassword = errors.New("incorrect email or password")
	NilPointer               = errors.New("nil pointer")
)
