package models

type RequestCreator struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}
type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RequestRegistration struct {
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
