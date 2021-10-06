package app

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"patreon/internal/app/sessions"
	"patreon/internal/app/store"
)

type ExpectedConnections struct {
	RedisPool     *redis.Pool
	SqlConnection *sql.DB
}

type DataStorage interface {
	Store() store.Store
	SessionManager() sessions.SessionsManager
}
