package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type UpdatePost struct {
	ID          int64  `json:"posts_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Awards      int64  `json:"type_awards"`
}

type CreatePost struct {
	ID          int64  `json:"posts_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Awards      int64  `json:"type_awards"`
	CreatorId   int64  `json:"creator_id"`
}

type Post struct {
	ID          int64     `json:"posts_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Awards      int64     `json:"type_awards"`
	Likes       int64     `json:"likes"`
	Cover       string    `json:"cover"`
	CreatorId   int64     `json:"creator_id"`
	Views       int64     `json:"views"`
	AddLike     bool      `json:"add_like"`
	Date        time.Time `json:"date"`
}

func (ps *UpdatePost) String() string {
	return fmt.Sprintf("{ID: %d, Title: %s, Likes: %d}", ps.ID,
		ps.Title, ps.Awards)
}

// Validate Errors:
//		EmptyTitle
//		InvalidAwardsId
// Important can return some other error
func (ps *UpdatePost) Validate() error {
	err := validation.Errors{
		"title":  validation.Validate(ps.Title, validation.Required),
		"awards": validation.Validate(ps.Awards, validation.Min(-1)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(postValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

func (ps *CreatePost) String() string {
	return fmt.Sprintf("{ID: %d, Title: %s, Likes: %d}", ps.ID,
		ps.Title, ps.Awards)
}

// Validate Errors:
//		EmptyTitle
//		InvalidCreatorId
//		InvalidAwardsId
// Important can return some other error
func (ps *CreatePost) Validate() error {
	err := validation.Errors{
		"title":   validation.Validate(ps.Title, validation.Required),
		"creator": validation.Validate(ps.CreatorId, validation.Min(0)),
		"awards":  validation.Validate(ps.Awards, validation.Min(-1)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(postValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

type DataType string

type PostData struct {
	ID     int64    `json:"data_id"`
	PostId int64    `json:"posts_id"`
	Data   string   `json:"data"`
	Type   DataType `json:"type"`
}

const (
	Music DataType = "music"
	Video DataType = "video"
	Files DataType = "files"
	Text  DataType = "text"
	Image DataType = "image"
)

// Validate Errors:
//		InvalidType
//		InvalidPostId
// Important can return some other error
func (ps *PostData) Validate() error {
	err := validation.Errors{
		"post": validation.Validate(ps.PostId, validation.Min(0)),
		"type": validation.Validate(ps.Type, validation.In(Music, Video, Files, Text, Image)),
	}.Filter()
	if err == nil {
		return nil
	}

	mapOfErr, knowError := parseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate creator")
	}

	if knowError = extractValidateError(postDataValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return err
}

type PostWithData struct {
	*Post
	Data []PostData
}
