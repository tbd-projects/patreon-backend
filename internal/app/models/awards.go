package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"strconv"
)

type Awards struct {
	ID          int64  `json:"awards_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	CreatorId   int64  `json:"creator_id"`
}

func (aw *Awards) String() string {
	return fmt.Sprintf("{ID: %s, Name: %s Price: %s}", strconv.Itoa(int(aw.ID)),
		aw.Name, strconv.Itoa(int(aw.Price)))
}

// Validate Errors:
//		EmptyName
//		IncorrectAwardsPrice
// Important can return some other error
func (aw *Awards) Validate() error {
	err := validation.Errors{
		"Name":    validation.Validate(aw.Name, validation.Required),
		"Price":    validation.Validate(aw.Price, validation.Min(0)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(awardsValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
