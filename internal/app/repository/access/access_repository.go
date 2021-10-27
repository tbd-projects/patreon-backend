package repository_access

import (
	"patreon/internal/app"

	"github.com/pkg/errors"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type RedisRepository struct {
	redisPool *redis.Pool
	log       *logrus.Logger
}

func NewRedisRepository(pool *redis.Pool, log *logrus.Logger) *RedisRepository {
	return &RedisRepository{
		redisPool: pool,
		log:       log,
	}
}

// Set Errors:
// 		app.GeneralError with Errors
// 			SetError
func (repo *RedisRepository) Set(key string, value string, timeExp int) error {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with key: %s value: %s",
				err.Error(), key, value)
		}
	}(con)
	res, err := redis.String(con.Do("SET", key, value, "EX", timeExp))
	if res != "OK" {
		return app.GeneralError{
			Err: errors.Wrapf(SetError,
				"error when try set with key: %s value: %s", key, value),
			ExternalErr: err,
		}
	}
	return nil
}

// Get Errors:
//		NotFound
// 		app.GeneralError with Errors
// 			InvalidStorageData
func (repo *RedisRepository) Get(key string) (string, error) {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with key: %s",
				err.Error(), key)
		}
	}(con)
	res, err := redis.String(con.Do("GET", key))
	if err == redis.ErrNil {
		return "", NotFound
	}
	if err != nil {
		return "", app.GeneralError{
			Err: errors.Wrapf(InvalidStorageData,
				"error when try get from AccessRepository with key: %s", key),
			ExternalErr: err,
		}
	}

	return res, nil
}

// Increment Errors:
// 		app.GeneralError with Errors
// 			InvalidStorageData
func (repo *RedisRepository) Increment(userIp string) (int64, error) {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with userIp: %s",
				err.Error(), userIp)
		}
	}(con)
	res, err := redis.Int64(con.Do("INCR", userIp))
	if err != nil {
		return -1, app.GeneralError{
			Err: InvalidStorageData,
			ExternalErr: errors.Wrapf(err,
				"error when try update userAccessCounter with userIp: %s", userIp),
		}
	}
	return res, nil
}
