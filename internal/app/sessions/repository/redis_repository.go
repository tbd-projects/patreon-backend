package repository

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"patreon/internal/app/sessions/models"
)

type RedisRepository struct {
	redisPool *redis.Pool
	log       *logrus.Logger
}

func CreateRedisRepository(pool *redis.Pool, log *logrus.Logger) *RedisRepository {
	return &RedisRepository{pool, log}
}

func (rep *RedisRepository) Set(session *models.Session) error {
	con := rep.redisPool.Get()
	defer func() {
		err := con.Close()
		if err != nil {
			rep.log.Errorf("Unsuccessful close connection to redis with error: %s, with session: %s",
				err.Error(), session)
		}
	}()

	result, err := redis.String(con.Do("SET", session.UniqID, session.UserID,
		"EX", session.Expiration))
	if result != "OK" {
		return errors.New("status not OK")
	}
	return err
}

func (rep *RedisRepository) GetUserId(uniqID string) (string, error) {
	con := rep.redisPool.Get()
	defer func() {
		err := con.Close()
		if err != nil {
			rep.log.Errorf("Unsuccessful close connection to redis with error: %s, with session id: %s",
				err.Error(), uniqID)
		}
	}()

	return redis.String(con.Do("GET", uniqID))
}

func (rep *RedisRepository) Del(session *models.Session) error {
	con := rep.redisPool.Get()
	defer func() {
		err := con.Close()
		if err != nil {
			rep.log.Errorf("Unsuccessful close connection to redis with error: %s, with session: %s",
				err.Error(), session)
		}
	}()

	_, err := redis.Int(con.Do("DEL", session.UniqID))
	return err
}
