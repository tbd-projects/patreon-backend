package models

import (
	"fmt"
	models_csrf "patreon/internal/app/csrf/models"
	"patreon/internal/app/models"
	"strconv"
)

type TokenResponse struct {
	Token models_csrf.Token `json:"token"`
}
type ErrResponse struct {
	Err string `json:"error"`
}

type UserResponse struct {
	ID int64 `json:"id"`
}

type AwardsResponse struct {
	ID int64 `json:"id"`
}

type ProfileResponse struct {
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ResponseCreator struct {
	models.Creator
}

type ResponseAwards struct {
	ID          int64  `json:"awards_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Color      	Color `json:"color,omitempty"`
}

func ToResponseCreator(cr models.Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func ToResponseAwards(aw models.Awards) ResponseAwards {
	return ResponseAwards{
		ID: aw.ID,
		Name: aw.Name,
		Price: aw.Price,
		Description: aw.Description,
		Color: NewColor(aw.Color),
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(int(u.ID)), u.Nickname)
}

type SubscriptionsUserResponse struct {
	Creators []int64 `json:"creator_id"`
}

func ToSubscriptionsUser(creators []int64) SubscriptionsUserResponse {
	return SubscriptionsUserResponse{
		Creators: creators,
	}
}
