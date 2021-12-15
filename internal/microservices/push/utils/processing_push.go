package utils

import "patreon/pkg/rabbit"

type SendMessager interface {
	SendMessage(users []int64, hsg interface{})
}

type ProcessingPush struct {
	session *rabbit.Session
	sendMsg SendMessager
}

func NewProcessingPush(session *rabbit.Session, sendMsg SendMessager) *ProcessingPush {
	return &ProcessingPush{
		session: session,
		sendMsg: sendMsg,
	}
}
