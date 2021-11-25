package utils

import "github.com/pkg/errors"

var (
	UnknownExtOfFileName = errors.New("Not found ext in file name")
	ConvertErr = errors.New("error of convert")
)
