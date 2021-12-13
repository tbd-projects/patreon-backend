package models_utilits

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ExtractorErrorByName func(string) error
type MapOfValidateError map[string]error
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
	encoded, bad := json.Marshal(err)
	if bad != nil {
		return nil, bad
	}

	bad = json.Unmarshal(encoded, &mapOfErr)
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
