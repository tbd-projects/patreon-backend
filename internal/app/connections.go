package app

import (
	"github.com/jmoiron/sqlx"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	SessionRedisPool *redis.Pool
	AccessRedisPool  *redis.Pool
	SqlConnection    *sqlx.DB
	PathFiles        string
}
