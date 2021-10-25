package models

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

type RequestAwards struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price"`
	Color       uint64 `json:"color.omitempty"`
}
