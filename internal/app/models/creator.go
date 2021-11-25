package models

import (
	"fmt"
	models_utilits "patreon/internal/app/utilits/models"
	"strconv"

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

type CreatorWithAwards struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	AwardsId    int64  `json:"awards_id"`
}

type CreatorSubscribe struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	AwardsId    int64  `json:"awards_id"`
}

func (cr *Creator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s Category: %s}", strconv.Itoa(int(cr.ID)), cr.Nickname, cr.Category)
}

// Validate Errors:
//		IncorrectCreatorNickname
//		IncorrectCreatorCategory
//		IncorrectCreatorDescription
// Important can return some other error
func (cr *Creator) Validate() error {
	err := validation.Errors{
		"nickname":    validation.Validate(cr.Nickname, validation.Required),
		"category":    validation.Validate(cr.Category, validation.Required),
		"description": validation.Validate(cr.Description, validation.Required),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = models_utilits.ExtractValidateError(creatorValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
