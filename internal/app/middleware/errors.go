package middleware

import "github.com/pkg/errors"

var (
	ForbiddenChangeCreator   = errors.New("for this user forbidden change creator")
	IncorrectCreatorForPost  = errors.New("this post not belongs this creators")
	IncorrectCreatorForAward = errors.New("this award not belongs this creators")
	InvalidParameters        = errors.New("invalid parameters")
	ContextError             = errors.New("can not get info from context")
	BDError                  = errors.New("can not do bd operation")
)
