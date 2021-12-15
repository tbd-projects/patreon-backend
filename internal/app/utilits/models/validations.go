package models_utilits

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mailru/easyjson"
)

//go:generate easyjson -disallow_unknown_fields validations.go

type ExtractorErrorByName func(string) error
type MapOfValidateError map[string]error

//easyjson:json
type MapOfUnmarshalError map[string]string

func RequiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}

func ParseErrorToMap(err error) (MapOfUnmarshalError, error) {
	var mapOfErr MapOfUnmarshalError
	encoded, bad := err.(validation.Errors).MarshalJSON()
	if bad != nil {
		return nil, bad
	}

	bad = easyjson.Unmarshal(encoded, &mapOfErr)
	if bad != nil {
		return nil, bad
	}
	return mapOfErr, nil
}

func ExtractValidateError(extractor ExtractorErrorByName, mapOfErr MapOfUnmarshalError) error {
	var knowError error = nil
	for key := range mapOfErr {
		if knowError = extractor(key); knowError != nil {
			break
		}
	}
	return knowError
}
