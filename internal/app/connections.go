package app

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	SessionGrpcConnection *grpc.ClientConn
	AccessRedisPool       *redis.Pool
	SqlConnection         *sql.DB
	PathFiles             string
}
