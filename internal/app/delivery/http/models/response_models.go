package models

import (
	"fmt"
	"patreon/internal/app/models"
	"strconv"
)

type BaseResponse struct {
	Code int    `json:"status"`
	Err  string `json:"error"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
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
