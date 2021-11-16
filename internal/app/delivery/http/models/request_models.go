package models

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"image/color"
	"patreon/internal/app/models"
	rep "patreon/internal/app/repository"
	models_utilits "patreon/internal/app/utilits/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	AddStatus    = "add"
	UpdateStatus = "update"
)

type RequestCreator struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RequestChangePassword struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}
type RequestChangeNickname struct {
	OldNickname string `json:"old"`
	NewNickname string `json:"new"`
}
type RequestRegistration struct {
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type Color struct {
	R uint8 `json:"red"`
	G uint8 `json:"green"`
	B uint8 `json:"blue"`
	A uint8 `json:"alpha"`
}

func NewColor(rgba color.RGBA) Color {
	return Color{
		R: rgba.R,
		G: rgba.G,
		B: rgba.B,
		A: rgba.A,
	}
}

type RequestAwards struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price"`
	Color       Color  `json:"color,omitempty"`
}

type RequestPosts struct {
	Title       string `json:"title,omitempty"`
	AwardsId    int64  `json:"awards_id,omitempty"`
	Description string `json:"description,omitempty"`
	IsDraft     bool   `json:"is_draft,omitempty"`
}

type RequestAttach struct {
	Type   string `json:"type"`
	Value  string `json:"value,omitempty"`
	Id     int64  `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type RequestAttaches struct {
	Attaches []RequestAttach `json:"attaches"`
}

type RequestText struct {
	Text string `json:"text"`
}

func (o *RequestPosts) UnmarshalJSON(text []byte) error {
	type options RequestPosts
	opts := options{
		AwardsId: rep.NoAwards,
	}

	if err := json.Unmarshal(text, &opts); err != nil {
		return err
	}

	*o = RequestPosts(opts)
	return nil
}

type SubscribeRequest struct {
	AwardName string `json:"award_name"`
}

func (req *SubscribeRequest) Validate() error {
	err := validation.Errors{
		"award_name": validation.Validate(req.AwardName, validation.Required, validation.Length(1, 0)),
	}.Filter()
	if err != nil {
		return AwardNameValidateError
	}
	return nil
}

func (req *RequestChangeNickname) Validate() error {
	err := validation.Errors{
		"old_nickname": validation.Validate(req.OldNickname, validation.Required,
			validation.Length(models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH)),
		"new_nickname": validation.Validate(req.NewNickname, validation.Required,
			validation.Length(models.MIN_NICKNAME_LENGTH, models.MAX_NICKNAME_LENGTH)),
	}.Filter()
	if err != nil {
		return NicknameValidateError
	}
	return nil
}

var (
	IncorrectType = errors.New(
		fmt.Sprintf("Not allow type, allowed type is: %s, %s, %s, %s, %s",
			models.Music, models.Video, models.Files, models.Text, models.Image))
	IncorrectIdAttach = errors.New("Not valid attach id")
	IncorrectStatus   = errors.New(fmt.Sprintf("Not allow status, allowed status is: %s, %s",
		AddStatus, UpdateStatus))
)

// requestAttachValidError Errors:
//		InvalidType
//		InvalidPostId
func requestAttachValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"type":   IncorrectType,
		"id":     IncorrectIdAttach,
		"status": IncorrectStatus,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}


func (req *RequestAttach) Validate() error {
	err := validation.Errors{
		"type": validation.Validate(req.Type, validation.In(models.Music, models.Video,
			models.Files, models.Text, models.Image)),
		"id":     validation.Validate(req.Id, validation.Min(1)),
		"status": validation.Validate(req.Status, validation.In(AddStatus, UpdateStatus)),
	}

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate request attach")
	}
	_, haveTypeError := mapOfErr["type"]
	_, haveIdError := mapOfErr["id"]
	_, haveStatusError := mapOfErr["status"]
	if !haveTypeError && haveIdError {
		if haveStatusError && models.DataType(req.Type) != models.Text {
			return IncorrectStatus
		}
		return nil
	}

	if haveStatusError && models.DataType(req.Type) != models.Text {
		if haveIdError {
			return IncorrectIdAttach
		}
		return nil
	}

	if knowError = models_utilits.ExtractValidateError(requestAttachValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return nil
}

func (req *RequestAttaches) Validate() error {
	for _, attach := range req.Attaches {
		if err := attach.Validate(); err != nil {
			return err
		}
	}
	return nil
}