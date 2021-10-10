package repository_creator

import "errors"

var (
	NotFound         = errors.New("user not found")
	UserAlreadyExist = errors.New("user already exist")
)
