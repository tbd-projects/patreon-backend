package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type Creator struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
}

// Validate Errors:
//		IncorrectCreatorNickname
//		IncorrectCreatorCategory
//		IncorrectCreatorCategoryDescription
// Important can return some other error
func (cr *Creator) Validate() error {
	err := validation.Errors{
		"id":          validation.Validate(cr.Nickname, validation.Required),
		"category":    validation.Validate(cr.Category, validation.Required),
		"description": validation.Validate(cr.Description, validation.Required),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if err != nil {
		return errors.Wrap(err, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(creatorValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
