package models

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
	Creator
}
