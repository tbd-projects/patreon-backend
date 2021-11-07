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
type OkResponse struct {
	Ok string `json:"OK"`
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

type ResponseInfo struct {
	models.Info
}

type ResponseCreatorWithAwards struct {
	models.CreatorWithAwards
}

type ResponseCreator struct {
	models.Creator
}

type ResponseCreatorSubscrube struct {
	models.CreatorSubscribe
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

func ToResponseInfo(info models.Info) ResponseInfo {
	return ResponseInfo{
		info,
	}
}

func ToResponseCreatorWithAwards(cr models.CreatorWithAwards) ResponseCreatorWithAwards {
	res := ResponseCreatorWithAwards{
		cr,
	}
	res.AwardsId = int64(math.Max(float64(res.AwardsId), 0))
	return res
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
	Creators []ResponseCreatorSubscrube `json:"creators"`
}

func ToSubscriptionsUser(creators []models.CreatorSubscribe) SubscriptionsUserResponse {
	var res []ResponseCreatorSubscrube
	for _, creator := range creators {
		res = append(res, ResponseCreatorSubscrube{
			models.CreatorSubscribe{
				ID:          creator.ID,
				Nickname:    creator.Nickname,
				Description: creator.Description,
				Category:    creator.Category,
				Cover:       creator.Cover,
				Avatar:      creator.Avatar,
				AwardsId:    creator.AwardsId,
			},
		})
	}
	return SubscriptionsUserResponse{
		Creators: res,
	}
}

type ResponseUser struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}
type SubscribersCreatorResponse struct {
	Users []ResponseUser `json:"users"`
}

func ToSubscribersCreatorResponse(users []models.User) SubscribersCreatorResponse {
	res := make([]ResponseUser, 0, len(users))
	for i, user := range users {
		res[i] = ResponseUser{
			ID:       user.ID,
			Login:    user.Login,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		}
	}
	return SubscribersCreatorResponse{
		Users: res,
	}
}

type ResponseLike struct {
	Likes int64 `json:"likes"`
}

type ResponseUserPayments struct {
	Payments []models.Payments `json:"payments"`
}

func ToResponseUserPayments(payments []models.Payments) ResponseUserPayments {
	res := make([]models.Payments, 0, len(payments))
	for _, payment := range payments {
		res = append(res, models.Payments{
			Amount:    payment.Amount,
			Date:      payment.Date,
			CreatorID: payment.CreatorID,
		})
	}
	return ResponseUserPayments{
		Payments: res,
	}
}
