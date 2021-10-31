package repository_postgresql

import (
	"github.com/pkg/errors"
)

var (
	NameAlreadyExist  = errors.New("name already exist")
	AwardNameNotFound = errors.New("creator have not this awardName")
)
