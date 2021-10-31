package models

import (
	"fmt"
	"math"
	models_csrf "patreon/internal/app/csrf/models"
	"patreon/internal/app/models"
	"strconv"
	"time"
)

type TokenResponse struct {
	Token models_csrf.Token `json:"token"`
}
type ErrResponse struct {
	Err string `json:"error"`
}

type IdResponse struct {
	ID int64 `json:"id"`
}

type ProfileResponse struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	HaveCreator bool   `json:"have_creator"`
}

type ResponseCreator struct {
	models.Creator
}

type ResponseAward struct {
	ID          int64  `json:"awards_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Color       Color  `json:"color,omitempty"`
	Cover       string `json:"cover"`
}

type ResponsePost struct {
	ID          int64     `json:"posts_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Awards      int64     `json:"type_awards,omitempty"`
	Likes       int64     `json:"likes"`
	Cover       string    `json:"cover"`
	AddLike     bool      `json:"add_like,omitempty"`
	Views       int64     `json:"views"`
	Date        time.Time `json:"date"`
}

type ResponsePostData struct {
	ID   int64  `json:"attach_id"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type ResponsePostWithData struct {
	Post ResponsePost       `json:"post"`
	Data []ResponsePostData `json:"attach"`
}
type ResponseBalance struct {
	ID      int64        `json:"user_id"`
	Balance models.Money `json:"balance"`
}

func ToRProfileResponse(us models.User) ProfileResponse {
	return ProfileResponse{
		ID:          us.ID,
		Login:       us.Login,
		Nickname:    us.Nickname,
		Avatar:      us.Avatar,
		HaveCreator: us.HaveCreator,
	}
}

func ToResponseCreator(cr models.Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func ToResponseAward(aw models.Award) ResponseAward {
	return ResponseAward{
		ID:          aw.ID,
		Name:        aw.Name,
		Price:       aw.Price,
		Description: aw.Description,
		Color:       NewColor(aw.Color),
		Cover:       aw.Cover,
	}
}

func ToResponsePost(ps models.Post) ResponsePost {
	return ResponsePost{
		ID:          ps.ID,
		Title:       ps.Title,
		Description: ps.Description,
		Date:        ps.Date,
		Likes:       ps.Likes,
		Awards:      int64(math.Max(float64(ps.Awards), 0)),
		Cover:       ps.Cover,
		AddLike:     ps.AddLike,
		Views:       ps.Views,
	}
}

func ToResponsePostWithData(ps models.PostWithData) ResponsePostWithData {
	res := ResponsePostWithData{Post: ToResponsePost(*ps.Post), Data: []ResponsePostData{}}
	for _, data := range ps.Data {
		res.Data = append(res.Data, ToResponsePostData(data))
	}
	return res
}

func ToResponsePostData(ps models.PostData) ResponsePostData {
	return ResponsePostData{
		ID:   ps.ID,
		Data: ps.Data,
		Type: string(ps.Type),
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

type SubscribersCreatorResponse struct {
	Users []int64 `json:"user_id"`
}

func ToSubscribersCreatorResponse(users []int64) SubscribersCreatorResponse {
	return SubscribersCreatorResponse{
		Users: users,
	}
}
