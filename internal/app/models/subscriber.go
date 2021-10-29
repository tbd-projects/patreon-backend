package models

type Subscriber struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"users_id"`
	CreatorID int64 `json:"creator_id"`
}
