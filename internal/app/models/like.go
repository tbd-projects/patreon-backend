package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type Like struct {
	ID     int64 `json:"likes_id"`
	Value  int8  `json:"value"`
	PostId int64 `json:"posts_id"`
	UserId int64 `json:"users_id"`
}

func (lk *Like) String() string {
	return fmt.Sprintf("{ID: %d, Value: %d PostID: %d UserID: %d}", lk.ID,
		lk.Value, lk.PostId, lk.UserId)
}

// Validate Errors:
//		EmptyName
//		IncorrectlkardsPrice
// Important can return some other error
func (lk *Like) Validate() error {
	err := validation.Errors{
		"value": validation.Validate(lk.Value, validation.By(func(val interface{}) error {
			value, ok := val.(*int8)
			if !ok || (*value != -1 && *value != 1) {
				return errors.New(fmt.Sprintf("Unvalid like value %v", value))
			}
			return nil
		})),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(likeValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
