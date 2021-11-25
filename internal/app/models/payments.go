package models

import "time"

type Payments struct {
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	CreatorID int64     `json:"creator_id,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
}

type UserPayments struct {
	Payments
	CreatorNickname    string `json:"creator_nickname"`
	CreatorCategory    string `json:"creator_category"`
	CreatorDescription string `json:"creator_description"`
}

type CreatorPayments struct {
	Payments
	UserNickname string `json:"user_nickname"`
}
