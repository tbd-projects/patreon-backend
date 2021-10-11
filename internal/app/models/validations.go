package models

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
)

type extractorErrorByName func(string)error
type mapOfValidateError map[string]error

func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}


func parseErrorToMap(err error) (map[string]error, error) {
	var mapOfErr map[string]error
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

func extractValidateError(extractor extractorErrorByName, mapOfErr mapOfValidateError) error {
	var knowError error = nil
	for key, _ := range mapOfErr {
		if knowError = extractor(key); knowError != nil {
			break
		}
	}
	return knowError
}
