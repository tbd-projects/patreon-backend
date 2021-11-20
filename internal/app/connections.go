package app

import (
	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	SessionGrpcConnection *grpc.ClientConn
	AccessRedisPool       *redis.Pool
	SqlConnection         *sqlx.DB
	PathFiles             string
}
