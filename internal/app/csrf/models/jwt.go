package models

import (
	"time"
)

type Token string

type TokenSources struct {
	UserId      int64
	SessionId   string
	ExpiredTime time.Time
}
