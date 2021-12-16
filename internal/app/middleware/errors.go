package middleware

import "github.com/pkg/errors"

var (
	ForbiddenChangeCreator   = errors.New("for this user forbidden change creator")
	IncorrectCreatorForPost  = errors.New("this post not belongs this creators")
	IncorrectCommentForPost  = errors.New("this comment not belongs this post")
	IncorrectCommentForUser  = errors.New("this comment not belongs this user")
	IncorrectAttachForPost   = errors.New("this attach not belongs this post")
	IncorrectCreatorForAward = errors.New("this award not belongs this creators")
	InvalidParameters        = errors.New("invalid parameters")
	BDError                  = errors.New("can not do bd operation")
	InternalError            = errors.New("server error")
)
