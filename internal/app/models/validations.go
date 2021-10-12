package models

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

type extractorErrorByName func(string) error
type mapOfValidateError map[string]error
type mapOfUnmarshalError map[string]string

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}

func parseErrorToMap(err error) (mapOfUnmarshalError, error) {
	var mapOfErr mapOfUnmarshalError
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

func extractValidateError(extractor extractorErrorByName, mapOfErr mapOfUnmarshalError) error {
	var knowError error = nil
	for key, _ := range mapOfErr {
		if knowError = extractor(key); knowError != nil {
			break
		}
	}
	return knowError
}
