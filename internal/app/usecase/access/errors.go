package usecase_access

import "errors"

var (
	NoAccess   = errors.New("exceeded the number of requests, repeat later")
	FirstQuery = errors.New("first query in timestamp")
)
