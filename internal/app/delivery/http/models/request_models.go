package models

import (
	"encoding/json"
	"image/color"
	repPosts "patreon/internal/app/repository/posts"
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
	Title       string `json:"title"`
	AwardsId    int64  `json:"awards_id,omitempty"`
	Description string `json:"description,omitempty"`
}

func (o *RequestPosts) UnmarshalJSON(text []byte) error {
	type options RequestPosts
	opts := options{
		AwardsId: repPosts.NoAwards,
	}

	if err := json.Unmarshal(text, &opts); err != nil {
		return err
	}

	*o = RequestPosts(opts)
	return nil
}

type RequestBalanceTop struct {
	Deposit int64 `json:"deposit"`
}
