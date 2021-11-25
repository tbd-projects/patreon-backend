package client

import (
	"bytes"
	"context"
	"io"
	grpc2 "patreon/internal/microservices/files/delivery/grpc"

	proto "patreon/internal/microservices/files/delivery/grpc/protobuf"
	repFiles "patreon/internal/microservices/files/files/repository/files"

	"google.golang.org/grpc"
)

type FileClient struct {
	fileClient proto.FileServiceClient
}

func NewFileServiceClient(con *grpc.ClientConn) *FileClient {
	client := proto.NewFileServiceClient(con)
	return &FileClient{
		fileClient: client,
	}
}

func (c *FileClient) SaveFile(ctx context.Context, file io.Reader, name repFiles.FileName, typeF repFiles.TypeFiles) (string, error) {
	fileToBytes := grpc2.StreamToByte(file)

	protoSaveArgs := &proto.SaveArgs{
		File: &proto.File{
			File: fileToBytes,
		},
		Name: string(name),
		Type: string(typeF),
	}
	res, err := c.fileClient.SaveFile(ctx, protoSaveArgs)
	if err != nil {
		return "", err
	}
	return res.Path, nil
}
func (c *FileClient) LoadFile(ctx context.Context, path string) (io.Reader, error) {
	protoLoadFile := &proto.FilePath{
		Path: path,
	}
	res, err := c.fileClient.LoadFile(ctx, protoLoadFile)
	if err != nil {
		return nil, err
	}
	resFile := bytes.NewReader(res.File)

	return resFile, err
}
