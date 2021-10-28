package models

import "time"

type Post struct {
	ID          int64     `json:"posts_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Likes       int64     `json:"likes"`
	Awards      int64     `json:"type_awards"`
	Cover       string    `json:"cover"`
	CreatorId   int64     `json:"creator_id"`
	Date        time.Time `json:"date"`
}

type PostData struct {
	ID     int64  `json:"data_id"`
	PostId int64 `json:"posts_id"`
	Data   string `json:"data"`
	Type   string `json:"type"`
}
