package utils

import (
	"bytes"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"patreon/internal/microservices/push"
	"patreon/internal/microservices/push/push/usecase"
	"patreon/pkg/rabbit"
)

type SendMessager interface {
	SendMessage(users []int64, hsg easyjson.Marshaler)
}

type ProcessingPush struct {
	session *rabbit.Session
	logger  *logrus.Entry
	sendMsg SendMessager
	usecase usecase.Usecase
	stop chan bool
}

func NewProcessingPush(logger *logrus.Entry, session *rabbit.Session, sendMsg SendMessager, usecase usecase.Usecase) *ProcessingPush {
	return &ProcessingPush{
		session: session,
		sendMsg: sendMsg,
		logger:  logger,
		usecase: usecase,
		stop : make(chan bool),
	}
}

func (pp *ProcessingPush) Stop() {
	pp.stop <- true
}

func (pp *ProcessingPush) initMsg(routerKey string) (<-chan amqp.Delivery, error) {
	ch := pp.session.GetChannel()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	if err = ch.QueueBind(q.Name, routerKey, pp.session.GetName(), false, nil); err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, err
}

func (pp *ProcessingPush) RunProcessPost() {
	msg, err := pp.initMsg(push.PostPush)
	if err != nil {
		pp.logger.Errorf("error init post query from msg with err: %s", err)
		return
	}
	pp.processPostMsg(msg)
}

func (pp *ProcessingPush) RunProcessComment() {
	msg, err := pp.initMsg(push.CommentPush)
	if err != nil {
		pp.logger.Errorf("error init comment query from msg with err: %s", err)
		return
	}
	pp.processCommentMsg(msg)
}

func (pp *ProcessingPush) RunProcessSub() {
	msg, err := pp.initMsg(push.NewSubPush)
	if err != nil {
		pp.logger.Errorf("error init sub query from msg with err: %s", err)
		return
	}
	pp.processSubMsg(msg)
}

func (pp *ProcessingPush) processPostMsg(msg <-chan amqp.Delivery) {
	for {
		var pushMsg amqp.Delivery
		select {
		case <-pp.stop:
			return
		case pushMsg = <-msg:
			break
		}

		post := &push.PostInfo{}
		reader := bytes.NewBuffer(pushMsg.Body)
		if err := easyjson.UnmarshalFromReader(reader, post); err != nil {
			pp.logger.Errorf("error decode info post from msg with err: %s", err)
			continue
		}

		users, sendPush, err := pp.usecase.PreparePostPush(post)
		if err != nil {
			pp.logger.Errorf("error prepare info post with err: %s", err)
			continue
		}
		pp.logger.Infof("Was send message about new post %v", pushMsg.Body)
		pp.sendMsg.SendMessage(users, PushResponse{Type: push.PostPush, Push: sendPush})
	}
}

func (pp *ProcessingPush) processCommentMsg(msg <-chan amqp.Delivery) {
	for {
		var pushMsg amqp.Delivery
		select {
		case <-pp.stop:
			return
		case pushMsg = <-msg:
			break
		}
		comment := &push.CommentInfo{}
		reader := bytes.NewBuffer(pushMsg.Body)
		if err := easyjson.UnmarshalFromReader(reader, comment); err != nil {
			pp.logger.Errorf("error decode info post from msg with err: %s", err)
			continue
		}

		users, sendPush, err := pp.usecase.PrepareCommentPush(comment)
		if err != nil {
			pp.logger.Errorf("error prepare info comment with err: %s", err)
			continue
		}
		pp.logger.Infof("Was send message about new comment %v", pushMsg.Body)
		pp.sendMsg.SendMessage(users, PushResponse{Type: push.CommentPush, Push: sendPush})
	}
}

func (pp *ProcessingPush) processSubMsg(msg <-chan amqp.Delivery) {
	for {
		var pushMsg amqp.Delivery
		select {
		case <-pp.stop:
			return
		case pushMsg = <-msg:
			break
		}

		subscriber := &push.SubInfo{}
		reader := bytes.NewBuffer(pushMsg.Body)
		if err := easyjson.UnmarshalFromReader(reader, subscriber); err != nil {
			pp.logger.Errorf("error decode info post from msg with err: %s", err)
			continue
		}

		users, sendPush, err := pp.usecase.PrepareSubPush(subscriber)
		if err != nil {
			pp.logger.Errorf("error prepare info sub with err: %s", err)
			continue
		}
		pp.logger.Infof("Was send message about new subscriber %v", pushMsg.Body)
		pp.sendMsg.SendMessage(users, PushResponse{Type: push.NewSubPush, Push: sendPush})
	}
}