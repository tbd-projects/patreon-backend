package push_client

type Pusher interface {
	NewPost(creatorId int64, postId int64, postTitle int64) error
	NewComment(authorId int64, postId int64) error
}