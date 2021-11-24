package client

import (
	"context"
	"io"
	repFiles "patreon/internal/microservices/files/files/repository/files"
)

//go:generate mockgen -destination=mocks/files_mock.go -package=mock_files . FileServiceClient

type FileServiceClient interface {
	SaveFile(ctx context.Context, file io.Reader, name repFiles.FileName, typeF repFiles.TypeFiles) (string, error)
	LoadFile(ctx context.Context, path string) (io.Reader, error)
}
