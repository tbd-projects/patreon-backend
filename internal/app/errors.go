package app

import (
	"fmt"
	"github.com/pkg/errors"
)


const InvalidInt = -1

var UnknownError = errors.New("gotten unspecified error")

type GeneralError struct {
	Err         error
	ExternalErr error
}

func (e GeneralError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.ExternalErr)
}
