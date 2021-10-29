package models

import (
	"fmt"
	"math"
	"patreon/internal/app/models"
	"strconv"
	"time"
)

type ErrResponse struct {
	Err string `json:"error"`
}

type IdResponse struct {
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

type ResponseAward struct {
	ID          int64  `json:"awards_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Color       Color  `json:"color,omitempty"`
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
	ID   int64  `json:"data_id"`
	Data string `json:"data"`
}

type ResponsePostWithData struct {
	ResponsePost
	Data []ResponsePostData `json:"data"`
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
	res := ResponsePostWithData{ResponsePost: ToResponsePost(*ps.Post), Data: []ResponsePostData{}}
	for _, data := range ps.Data {
		res.Data = append(res.Data, ToResponsePostData(data))
	}
	return res
}

func ToResponsePostData(ps models.PostData) ResponsePostData {
	return ResponsePostData{
		ID:   ps.ID,
		Data: ps.Data,
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(int(u.ID)), u.Nickname)
}
