package models

import (
	"fmt"
	"patreon/internal/app/models"
	"strconv"
)

type ErrResponse struct {
	Err string `json:"error"`
}

type UserResponse struct {
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

func ToResponseCreator(cr models.Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(int(u.ID)), u.Nickname)
}
