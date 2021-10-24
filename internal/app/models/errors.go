package models

import (
	"errors"
)

var (
	EmptyPassword                       = errors.New("empty password")
	IncorrectEmailOrPassword            = errors.New("invalid email or password")
	IncorrectCreatorNickname            = errors.New("incorrect creator nickname")
	EmptyName                           = errors.New("empty name")
	IncorrectAwardsPrice                = errors.New("incorrect awards price")
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

// creatorValidError Errors:
//		EmptyName
//		IncorrectAwardsPrice
func awardsValidError() extractorErrorByName {
	validMap := mapOfValidateError{
		"name":  EmptyName,
		"price": IncorrectAwardsPrice,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}
