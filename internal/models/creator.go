package models

import (
	"fmt"
	"strconv"
)

type Creator struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
}

type ResponseCreator struct {
	Creator
}

func ToResponseCreator(cr Creator) ResponseCreator {
	return ResponseCreator{
		cr,
	}
}

func (u *ResponseCreator) String() string {
	return fmt.Sprintf("{ID: %s, Nickname: %s}", strconv.Itoa(u.ID), u.Nickname)
}
