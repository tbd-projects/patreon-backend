package app

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	//"patreon/internal/app/store"
)

type ExpectedConnections struct {
	RedisPool     *redis.Pool
	SqlConnection *sql.DB
}