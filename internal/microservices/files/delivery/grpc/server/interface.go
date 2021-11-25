package server

import (
	"context"
	proto "patreon/internal/microservices/files/delivery/grpc/protobuf"
)

type FileServiceServer interface {
	SaveFile(ctx context.Context, args *proto.SaveArgs) (*proto.FilePath, error)
	LoadFile(ctx context.Context, path *proto.FilePath) (*proto.File, error)
}
