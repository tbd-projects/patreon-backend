package push_client

import (
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	models "patreon/internal/microservices/push"
	"patreon/pkg/rabbit"
	"time"
)

type PushSender struct {
	session *rabbit.Session
}

func NewPushSender(session *rabbit.Session) *PushSender {
	return &PushSender{
		session: session,
	}
}

func (ph *PushSender) push(routeName string, msg easyjson.Marshaler) error {
	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}
	var err error
	publish.Body, err = easyjson.Marshal(msg)
	if err != nil {
		return err
	}
	ch := ph.session.GetChannel()

	err = ch.Publish(
		ph.session.GetName(),
		routeName,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewPost(creatorId int64, postId int64, postTitle string) error {
	return ph.push(models.PostPush, &models.PostInfo{
		CreatorId: creatorId,
		PostId:    postId,
		PostTitle: postTitle,
		Date:      time.Now(),
	})
}

func (ph *PushSender) NewComment(commentId int64, authorId int64, postId int64) error {
	return ph.push(models.CommentPush, &models.CommentInfo{
		CommentId: commentId,
		AuthorId:  authorId,
		PostId:    postId,
		Date:      time.Now(),
	})
}

func (ph *PushSender) ApplyPayments(token string) error {
	return ph.push(models.PaymentPush, &models.PaymentApply{
		Token: token,
		Date:  time.Now(),
	})
}