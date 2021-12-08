package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"patreon/internal/app/utilits/models"
	"strconv"
)

type Comment struct {
	ID        int64  `json:"comment_id"`
	Body      string `json:"body"`
	AsCreator bool   `json:"as_creator,omitempty"`
	AuthorId  int64  `json:"author_id"`
	PostId    int64  `json:"post_id.omitempty"`
}

type PostComment struct {
	Comment
	AuthorNickname string `json:"author_nickname"`
	AuthorAvatar   string `json:"author_avatar"`
}

type UserComment struct {
	Comment
	PostName  string `json:"post_name"`
	PostCover string `json:"post_cover"`
}

func (cm *Comment) String() string {
	return fmt.Sprintf("{ID: %s, Body: %s postId: %s authorId %s}", strconv.Itoa(int(cm.ID)),
		cm.Body, strconv.Itoa(int(cm.PostId)), strconv.Itoa(int(cm.AuthorId)))
}

// Validate Errors:
//		InvalidUserId
//		InvalidPostId
// Important can return some other error
func (cm *Comment) Validate() error {
	err := validation.Errors{
		"author_id": validation.Validate(cm.AuthorId, validation.Min(1)),
		"post_id":   validation.Validate(cm.PostId, validation.Min(1)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = models_utilits.ExtractValidateError(commentsValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}
