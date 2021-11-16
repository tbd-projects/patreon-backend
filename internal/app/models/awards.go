package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"image/color"
	models_utilits "patreon/internal/app/utilits/models"
	"strconv"
)

type Award struct {
	ID          int64      `json:"awards_id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Price       int64      `json:"price,omitempty"`
	CreatorId   int64      `json:"creator_id"`
	Color       color.RGBA `json:"color.omitempty"`
	ChildAward  int64      `json:"child_award"`
	Cover       string     `json:"cover"`
}

func (aw *Award) String() string {
	return fmt.Sprintf("{ID: %s, Name: %s Price: %s}", strconv.Itoa(int(aw.ID)),
		aw.Name, strconv.Itoa(int(aw.Price)))
}

// Validate Errors:
//		EmptyName
//		IncorrectAwardsPrice
// Important can return some other error
func (aw *Award) Validate() error {
	err := validation.Errors{
		"name":  validation.Validate(aw.Name, validation.Required),
		"price": validation.Validate(aw.Price, validation.Min(0)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = models_utilits.ExtractValidateError(awardsValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
