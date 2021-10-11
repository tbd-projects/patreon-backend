package models

import (
	"errors"
)

var (
	EmptyPassword                       = errors.New("empty password")
	IncorrectEmailOrPassword            = errors.New("invalid email or password")
	IncorrectCreatorNickname            = errors.New("incorrect creator nickname")
	IncorrectCreatorCategory            = errors.New("incorrect creator category")
	IncorrectCreatorCategoryDescription = errors.New("incorrect creator category description")
)

// userValidError Errors:
//		IncorrectEmailOrPassword
func userValidError() extractorErrorByName {
	validMap := mapOfValidateError{
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

// creatorValidError Errors:
//		IncorrectCreatorNickname
//		IncorrectCreatorCategory
//		IncorrectCreatorCategoryDescription
func creatorValidError() extractorErrorByName {
	validMap := mapOfValidateError{
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
