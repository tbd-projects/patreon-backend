package usercase_user

import (
	"errors"
)

var (
	UserExist                = errors.New("user already exist")
	NicknameExists           = errors.New("this nickname already exist")
	BadEncrypt               = errors.New("unsuccessful encrypt user")
	IncorrectEmailOrPassword = errors.New("incorrect email or password")
	OldPasswordEqualNew      = errors.New("the new password must be different from the old one")
	IncorrectNewPassword     = errors.New("new password not valid")
)
