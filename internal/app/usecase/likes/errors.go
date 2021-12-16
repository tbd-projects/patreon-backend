package usecase_likes

import "github.com/pkg/errors"

var (
	IncorrectDelLike = errors.New("user try del likes that he not add")
	IncorrectAddLike = errors.New("user try add likes that already add")
)
