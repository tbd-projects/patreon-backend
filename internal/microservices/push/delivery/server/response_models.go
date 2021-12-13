package push_server

type PostResponse struct {
	PostId          int64  `json:"post_id"`
	CreatorId       int64  `json:"creator_id"`
	CreatorNickname string `json:"creator_nickname"`
	PostTitle       string `json:"post_title"`
	CreatorAvatar   string `json:"creator_avatar"`
}
