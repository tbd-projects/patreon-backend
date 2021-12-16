package push_client

type Pusher interface {
	NewPost(creatorId int64, postId int64, postTitle string) error
	NewComment(commentId int64, authorId int64, postId int64) error
	NewSubscriber(subscriberId int64, awardsId int64, creatorId int64) error
}
