package push_client

import (
	"bytes"
	"encoding/json"
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

func (ph *PushSender) NewPost(creatorId int64, postId int64, postTitle string) error {
	push := &models.PostInfo{
		CreatorId: creatorId,
		PostId:    postId,
		PostTitle: postTitle,
		Date:      time.Now(),
	}

	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}

	body := bytes.NewBuffer(publish.Body)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(push); err != nil {
		return err
	}
	ch, err := ph.session.GetChannel()
	if err != nil {
		return err
	}

	err = ch.Publish(
		ph.session.GetName(),
		models.PostPush,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewComment(commentId int64, authorId int64, postId int64) error {
	push := &models.CommentInfo{
		AuthorId: authorId,
		PostId:   postId,
		Date:     time.Now(),
	}

	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}

	body := bytes.NewBuffer(publish.Body)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(push); err != nil {
		return err
	}
	ch, err := ph.session.GetChannel()
	if err != nil {
		return err
	}

	err = ch.Publish(
		ph.session.GetName(),
		models.CommentPush,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewSubscriber(subscriberId int64, awardsId int64, creatorId int64) error {
	push := &models.SubInfo{
		UserId:     subscriberId,
		CreatorId: creatorId,
		AwardsId:   awardsId,
		Date:       time.Now(),
	}

	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}

	body := bytes.NewBuffer(publish.Body)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(push); err != nil {
		return err
	}
	ch, err := ph.session.GetChannel()
	if err != nil {
		return err
	}

	err = ch.Publish(
		ph.session.GetName(),
		models.CommentPush,
		false,
		false,
		publish,
	)

	return err
}
