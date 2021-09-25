package store

import "errors"

var (
	NotFound                 = errors.New("user not found")
	UserAlreadyExist         = errors.New("user already exist")
	IncorrectEmailOrPassword = errors.New("incorrect email or password")
	// @todo убрать ошибки в handlers
	GetProfileFail   = errors.New("can not get profile from db")
	DeleteCookieFail = errors.New("can not delete cookie from session store")
	InvalidBody      = errors.New("invalid body in request")
)
