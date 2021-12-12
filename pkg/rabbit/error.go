package rabbit

import "github.com/pkg/errors"

var (
	ErrAlreadyClosed = errors.New("already closed: not connected to the server")
	ErrGetUninnitChanel = errors.New("failed to get channel: not connected")
)