package server

import (
	"context"
	proto "patreon/internal/microservices/auth/delivery/grpc/protobuf"
)

type AuthCheckerServer interface {
	Check(context.Context, *proto.SessionID) (*proto.Result, error)
	Create(context.Context, *proto.UserID) (*proto.Result, error)
	Delete(context.Context, *proto.SessionID) (*proto.Nothing, error)
}
