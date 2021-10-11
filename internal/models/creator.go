package models

import (
	"encoding/json"
	"fmt"
	"patreon/internal/app"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Creator struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
}

func (cr *Creator) Validate() error {
	err := validation.Errors{
		"id":          validation.Validate(cr.Nickname, validation.Required),
		"category":    validation.Validate(cr.Category, validation.Required),
		"description": validation.Validate(cr.Description, validation.Required),
	}.Filter()
	if err == nil {
		return err
	}
	var mapOfErr map[string]error
	encoded, bad := json.Marshal(err)
	if bad != nil {
		return app.GeneralError{
			Err:         InternalError,
			ExternalErr: bad,
		}
	}
	bad = json.Unmarshal(encoded, &mapOfErr)
	if bad != nil {
		return app.GeneralError{
			Err:         InternalError,
			ExternalErr: bad,
		}
	}
	for key, _ := range mapOfErr {
		return creatorValidError()(key)
	}
	return InternalError
}
func ToResponseCreator(cr Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(int(u.ID)), u.Nickname)
}
