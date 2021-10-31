package app

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	SessionRedisPool *redis.Pool
	AccessRedisPool  *redis.Pool
	SqlConnection    *sql.DB
	PathFiles        string
}
