package models

import (
	"errors"
)

func userValidError() func(string) error {
	validMap := map[string]error{
		"login":    IncorrectEmailOrPassword,
		"password": IncorrectEmailOrPassword,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}
func creatorValidError() func(string) error {
	validMap := map[string]error{
		"nickname":    IncorrectCreatorNickname,
		"category":    IncorrectCreatorCategory,
		"description": IncorrectCreatorCategoryDescription,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

var (
	EmptyPassword                       = errors.New("empty password")
	IncorrectEmailOrPassword            = errors.New("invalid email or password")
	IncorrectCreatorNickname            = errors.New("incorrect creator nickname")
	IncorrectCreatorCategory            = errors.New("incorrect creator category")
	IncorrectCreatorCategoryDescription = errors.New("incorrect creator category description")
	InternalError                       = errors.New("internal error")
)
