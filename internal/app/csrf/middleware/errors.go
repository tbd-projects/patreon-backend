package middleware

import "errors"

var (
	InvalidToken = errors.New("invalid csrf token")
)
