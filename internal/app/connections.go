package app

import (
	"github.com/jmoiron/sqlx"
	"patreon/pkg/rabbit"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	SessionGrpcConnection *grpc.ClientConn
	FilesGrpcConnection   *grpc.ClientConn
	AccessRedisPool       *redis.Pool
	SqlConnection         *sqlx.DB
	PathFiles             string
	RabbitSession         *rabbit.Session
}
