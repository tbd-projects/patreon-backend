package repository_access

import (
	"patreon/internal/app"
	"patreon/internal/app/models"

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
func (repo *RedisRepository) Set(userIp string, access models.AccessCounter, timeExp int) error {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with userIp: %s",
				err.Error(), userIp)
		}
	}(con)
	res, err := redis.String(con.Do("SET", userIp, access.Counter, "EX", timeExp))
	if res != "OK" {
		return app.GeneralError{
			Err: errors.Wrapf(SetError,
				"error when try set userAccessCounter with userIp: %s", userIp),
			ExternalErr: err,
		}
	}
	return nil
}

// Get Errors:
//		NotFound
// 		app.GeneralError with Errors
// 			InvalidStorageData
func (repo *RedisRepository) Get(userIp string) (int64, error) {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with userIp: %s",
				err.Error(), userIp)
		}
	}(con)
	res, err := redis.Int64(con.Do("GET", userIp))
	if err == redis.ErrNil {
		return -1, NotFound
	}
	if err != nil {
		return -1, app.GeneralError{
			Err: errors.Wrapf(InvalidStorageData,
				"error when try get userAccessCounter with userIp: %s", userIp),
			ExternalErr: err,
		}
	}

	return res, nil
}

// Update Errors:
// 		app.GeneralError with Errors
// 			InvalidStorageData
func (repo *RedisRepository) Update(userIp string) (int64, error) {
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

// AddToBlackList Errors:
// 		app.GeneralError with Errors
// 			SetError
func (repo *RedisRepository) AddToBlackList(key string, userIp string, timeLimit int) error {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with key: %s",
				err.Error(), key)
		}
	}(con)
	blackListKey := key + userIp
	res, err := redis.String(con.Do("SET", blackListKey, userIp, "EX", timeLimit))
	if res != "OK" {
		return app.GeneralError{
			Err: SetError,
			ExternalErr: errors.Wrapf(err,
				"error when try add user with ip: %v in blackList", userIp),
		}
	}
	return nil
}

// CheckBlackList Errors:
//		NotFound
// 		app.GeneralError with Errors
// 			InvalidStorageData
func (repo *RedisRepository) CheckBlackList(key string, userIp string) (bool, error) {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with key: %s",
				err.Error(), key)
		}
	}(con)

	blackListKey := key + userIp
	_, err := redis.String(con.Do("GET", blackListKey))
	if err != nil {
		if err == redis.ErrNil {
			return false, NotFound
		} else {
			return true, app.GeneralError{
				Err: InvalidStorageData,
				ExternalErr: errors.Wrapf(err,
					"error when try get user from black list with ip: %v in blackList", userIp),
			}
		}
	}

	return true, nil
}
