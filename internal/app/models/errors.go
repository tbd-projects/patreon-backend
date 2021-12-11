package models

import (
	"errors"
	"fmt"
	models_utilits "patreon/internal/app/utilits/models"
)

var (
	EmptyPassword            = errors.New("empty password")
	IncorrectEmailOrPassword = errors.New("invalid email or password")
	IncorrectNickname        = errors.New(fmt.Sprintf("invalid nickname in body len must be from %v to %v",
		MIN_NICKNAME_LENGTH, MAX_NICKNAME_LENGTH))
	IncorrectCreatorNickname    = errors.New("incorrect creator nickname")
	EmptyName                   = errors.New("empty name")
	IncorrectAwardsPrice        = errors.New("incorrect awards price")
	IncorrectCreatorCategory    = errors.New("incorrect creator category")
	IncorrectCreatorDescription = errors.New("incorrect creator category description")
	InvalidLikeValue            = errors.New("like contain not valid value")
	EmptyTitle                  = errors.New("empty title")
	InvalidCreatorId            = errors.New("not positive creator id")
	InvalidAwardsId             = errors.New("not positive awards id")
	InvalidPostId               = errors.New("not positive posts id")
	InvalidUserId               = errors.New("not positive user id")
	InvalidType                 = errors.New("not positive data type")
)

// userValidError Errors:
//		IncorrectEmailOrPassword
//		IncorrectNickname
func userValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"login":    IncorrectEmailOrPassword,
		"password": IncorrectEmailOrPassword,
		"nickname": IncorrectNickname,
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
//		IncorrectCreatorDescription
func creatorValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"nickname":    IncorrectCreatorNickname,
		"category":    IncorrectCreatorCategory,
		"description": IncorrectCreatorDescription,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// awardsValidError Errors:
//		EmptyName
//		IncorrectAwardsPrice
func awardsValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
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

// awardsValidError Errors:
//		InvalidUserId
//		InvalidPostId
func commentsValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"author_id": InvalidUserId,
		"post_id":   InvalidPostId,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// likeValidError Errors:
//		EmptyName
//		IncorrectAwardsPrice
func likeValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"value": InvalidLikeValue,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// postValidError Errors:
//		EmptyTitle
//		InvalidCreatorId
//		InvalidAwardsId
func postValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"title":   EmptyTitle,
		"creator": InvalidCreatorId,
		"awards":  InvalidAwardsId,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// postValidError Errors:
//		InvalidType
//		InvalidPostId
func attachWithoutLevelValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"post": InvalidPostId,
		"type": InvalidType,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}
