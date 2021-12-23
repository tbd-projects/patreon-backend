package server

import (
	"context"
	proto "patreon/internal/microservices/smtp/delivery/grpc/protobuf"
)

type SmtpServiceServer interface {
	Send(ctx context.Context, message *proto.Message) error
	Stop(ctx context.Context) error
}
