package push_client

type Pusher interface {
	NewPost(creatorId int64, postId int64, postTitle string) error
	ApplyPayments(token string) error
	NewComment(commentId int64, authorId int64, postId int64) error
}
