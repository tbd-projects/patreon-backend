package app

import "fmt"

type GeneralError struct {
	Err         error
	ExternalErr error
}

func (e GeneralError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.ExternalErr)
}
