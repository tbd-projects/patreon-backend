package models

import (
	"encoding/json"
	"image/color"
	"patreon/internal/app/models"
	rep "patreon/internal/app/repository"

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
