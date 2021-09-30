package repository

import (
	"patreon/internal/app/sessions/models"

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

func (repo *RedisRepository) Set(session *models.Session) error {
	con := repo.redisPool.Get()
	defer func(con redis.Conn) {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with session: %s",
				err.Error(), session)
		}
	}(con)

	res, err := redis.String(con.Do("SET", session.UniqID, session.UserID,
		"EX", session.Expiration))
	if res != "OK" {
		return err
	}
	return nil
}
func (repo *RedisRepository) GetUserId(uniqID string) (string, error) {
	con := repo.redisPool.Get()
	defer func() {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with session id: %s",
				err.Error(), uniqID)
		}
	}()

	return redis.String(con.Do("GET", uniqID))
}

func (repo *RedisRepository) Del(session *models.Session) error {
	con := repo.redisPool.Get()
	defer func() {
		err := con.Close()
		if err != nil {
			repo.log.Errorf("Unsuccessful close connection to redis with error: %s, with session: %s",
				err.Error(), session)
		}
	}()

	_, err := redis.Int(con.Do("DEL", session.UniqID))
	return err
}
