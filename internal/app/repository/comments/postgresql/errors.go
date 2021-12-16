package repository_postgresql

import (
	"github.com/pkg/errors"
)

var (
	CommentAlreadyExist = errors.New("comment already exist")
)
