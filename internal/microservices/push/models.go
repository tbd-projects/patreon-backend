package push

import "time"

const (
	CommentPush = "Comment"
	PostPush    = string("Push")
)

type PostInfo struct {
	CreatorId int64     `json:"creator_id"`
	PostId    int64     `json:"post_id"`
	PostTitle int64     `json:"post_title"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
}

type CommentInfo struct {
	AuthorId int64     `json:"author_id"`
	PostId   int64     `json:"post_id"`
	Date     time.Time `json:"date"`
	Message  string    `json:"message"`
}
