package models

import "time"

type Payments struct {
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	CreatorID int64     `json:"creator_id"`
	UserID    int64     `json:"user_id"`
}
