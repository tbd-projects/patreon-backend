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
	CreatorId      int64  `json:"creator_id"`
	CommentId      int64  `json:"comment_id"`
	PostId         int64  `json:"post_id"`
	AuthorId       int64  `json:"author_id"`
	AuthorNickname string `json:"author_nickname"`
	AuthorAvatar   string `json:"author_avatar"`
	PostTitle      string `json:"post_title"`
}

//easyjson:json
type PaymentApplyPush struct {
	CreatorId       int64  `json:"creator_id"`
	CreatorNickname string `json:"creator_nickname"`
	CreatorAvatar   string `json:"creator_avatar"`
	AwardsId        int64  `json:"awards_id"`
	AwardsName      string `json:"awards_name"`
}
