package http_models

import (
	"encoding/json"
	"github.com/pkg/errors"
	"image/color"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/models"
	rep "patreon/internal/app/repository"
	models_utilits "patreon/internal/app/utilits/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RequestCreator struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RequestComment struct {
	Body      string `json:"body"`
	AsCreator bool   `json:"as_creator,omitempty"`
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
	Type   models.DataType `json:"type"`
	Value  string          `json:"value,omitempty"`
	Id     int64           `json:"id,omitempty"`
	Status string          `json:"status,omitempty"`
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

// requestAttachValidError Errors:
//		handler_errors.IncorrectType
//		handler_errors.IncorrectIdAttach
//      handler_errors.IncorrectStatus
func requestAttachValidError() models_utilits.ExtractorErrorByName {
	validMap := models_utilits.MapOfValidateError{
		"type":   handler_errors.IncorrectType,
		"id":     handler_errors.IncorrectIdAttach,
		"status": handler_errors.IncorrectStatus,
	}
	return func(key string) error {
		if val, ok := validMap[key]; ok {
			return val
		}
		return nil
	}
}

// Validate Errors:
//		handler_errors.IncorrectType
//		handler_errors.IncorrectIdAttach
//      handler_errors.IncorrectStatus
// can return not specify error
func (req *RequestAttach) Validate() error {
	err := validation.Errors{
		"type": validation.Validate(req.Type, validation.In(models.Music, models.Video,
			models.Files, models.Text, models.Image)),
		"id":     validation.Validate(req.Id, validation.Min(1)),
		"status": validation.Validate(req.Status, validation.In(handlers.AddStatus, handlers.UpdateStatus)),
	}.Filter()

	mapOfErr, knowError := models_utilits.ParseErrorToMap(err)
	if knowError != nil {
		return errors.Wrap(knowError, "failed error getting in validate request attach")
	}
	_, haveTypeError := mapOfErr["type"]
	_, haveIdError := mapOfErr["id"]
	_, haveStatusError := mapOfErr["status"]
	if !haveTypeError && haveIdError {
		if haveStatusError && req.Type != models.Text {
			return handler_errors.IncorrectStatus
		}
		return nil
	}

	if haveStatusError && req.Type != models.Text {
		if haveIdError {
			return handler_errors.IncorrectIdAttach
		}
		return nil
	}

	if knowError = models_utilits.ExtractValidateError(requestAttachValidError(), mapOfErr); knowError != nil {
		return knowError
	}

	return nil
}

// Validate Errors:
//		handler_errors.IncorrectType
//		handler_errors.IncorrectIdAttach
//      handler_errors.IncorrectStatus
// can return not specify error
func (req *RequestAttaches) Validate() error {
	for _, attach := range req.Attaches {
		if err := attach.Validate(); err != nil {
			return err
		}
	}
	return nil
}
