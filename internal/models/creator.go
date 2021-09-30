package models

import (
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Creator struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
}

func (cr *Creator) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.ID, validation.Required),
		validation.Field(&cr.Nickname, validation.Required),
		validation.Field(&cr.Category, validation.Required),
		validation.Field(&cr.Description, validation.Required),
	)
}
func ToResponseCreator(cr Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(u.ID), u.Nickname)
}
