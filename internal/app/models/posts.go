package models

import "time"

type Posts struct {
	ID          int64     `json:"posts_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Likes       int64     `json:"likes"`
	Awards      int64     `json:"type_awards"`
	Cover       int64     `json:"cover"`
	CreatorId   int64     `json:"creator_id"`
	Date        time.Time `json:"date"`
}

type PostsData struct {
	ID       int64  `json:"data_id"`
	PostId   string `json:"posts_id"`
	DataPath string `json:"data_path"`
	Type     int64  `json:"type"`
}
