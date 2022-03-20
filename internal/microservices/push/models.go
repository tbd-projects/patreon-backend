package push

import "time"

//go:generate easyjson -all -disallow_unknown_fields models.go

const (
	CommentPush = "Comment"
	PaymentPush = "Payment"
	PostPush    = "Post"
)

//easyjson:json
type PostInfo struct {
	CreatorId int64     `json:"creator_id"`
	PostId    int64     `json:"post_id"`
	PostTitle string    `json:"post_title"`
	Date      time.Time `json:"date"`
}

//easyjson:json
type CommentInfo struct {
	CommentId int64     `json:"comment_id"`
	AuthorId  int64     `json:"author_id"`
	PostId    int64     `json:"post_id"`
	Date      time.Time `json:"date"`
}

//easyjson:json
type PaymentApply struct {
	Token string    `json:"token"`
	Date  time.Time `json:"date"`
}
