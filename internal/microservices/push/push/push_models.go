package push_models

//go:generate easyjson -all -disallow_unknown_fields push_models.go

//easyjson:json
type PostPush struct {
	PostId          int64  `json:"post_id"`
	CreatorId       int64  `json:"creator_id"`
	CreatorNickname string `json:"creator_nickname"`
	PostTitle       string `json:"post_title"`
	CreatorAvatar   string `json:"creator_avatar"`
}

//easyjson:json
type CommentPush struct {
	CommentId      int64  `json:"comment_id"`
	PostId         int64  `json:"post_id"`
	AuthorId       int64  `json:"author_id"`
	AuthorNickname string `json:"author_nickname"`
	AuthorAvatar   string `json:"author_avatar"`
	PostTitle      string `json:"post_title"`
}

//easyjson:json
type SubPush struct {
	AwardsId     int64  `json:"awards_id"`
	UserId       int64  `json:"user_id"`
	AwardsName   string `json:"awards_name"`
	AwardsPrice  int64  `json:"awards_price"`
	UserNickname string `json:"user_nickname"`
	UserAvatar   string `json:"user_avatar"`
}
