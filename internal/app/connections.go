package app

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

type ExpectedConnections struct {
	RedisPool     *redis.Pool
	SqlConnection *sql.DB
}
