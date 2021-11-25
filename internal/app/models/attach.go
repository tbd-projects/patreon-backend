package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	models_utilits "patreon/internal/app/utilits/models"
)

var (
	IncorrectType     = errors.New("incorrect attach type")
	IncorrectAttachId = errors.New("incorrect attach id")
	IncorrectLevel    = errors.New("incorrect attach level")
)

type Attach struct {
	Id    int64    `json:"id"`
	Value string   `json:"value"`
	Type  DataType `json:"type"`
	Level int64    `json:"level"`
}

// AttachValidError Errors:
//		IncorrectType
//		IncorrectAttachId
//      IncorrectLevel
func AttachValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"type":  IncorrectType,
		"id":    IncorrectAttachId,
		"level": IncorrectLevel,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// Validate Errors:
//		IncorrectType
//		IncorrectAttachId
//      IncorrectLevel
// can return not specify error
func (att *Attach) Validate() error {
	err := validation.Errors{
		"type": validation.Validate(att.Type, validation.In(Music, Video,
			Files, Text, Image)),
		"id":    validation.Validate(att.Id, validation.Min(0)),
		"level": validation.Validate(att.Level, validation.Min(1)),
	}.Filter()

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate attach")
	}

	if knowError = models_utilits.ExtractValidateError(AttachValidError(), mapOfErr); knowError != nil {
		return errors.Wrap(knowError,
			fmt.Sprintf("error with attach id: %d, type: %s, level: %d", att.Id, att.Type, att.Level))
	}

	return nil
}
