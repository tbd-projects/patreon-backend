package models

type Posts struct {
	ID          int64  `json:"posts_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DataPath    string `json:"data_path"`
	Type        int64  `json:"type"`
	Likes       int64  `json:"likes"`
	Awards      int64  `json:"type_awards"`
	CreatorId   int64  `json:"creator_id"`
}
