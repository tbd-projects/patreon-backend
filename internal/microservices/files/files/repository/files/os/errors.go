package repository_os

import "github.com/pkg/errors"

var (
	ErrorCreate   = errors.New("some error create dir of files: ")
	ErrorCopyFile = errors.New("some error when data are being copied: ")
	ErrorOpenFile = errors.New("some error when data are being opened: ")
)
