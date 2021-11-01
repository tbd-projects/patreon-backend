package repository_access

import (
	"bytes"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SuiteTestRepository struct {
	suite.Suite
	redisServer      *miniredis.Miniredis
	accessRepository *RedisRepository
	log              *logrus.Logger
	output           string
}

func (s *SuiteTestRepository) SetupSuite() {
	s.log = logrus.New()
	s.log.SetLevel(logrus.FatalLevel)
	s.output = ""
	s.log.SetOutput(bytes.NewBufferString(s.output))

	var err error
	s.redisServer, err = miniredis.Run()
	require.NoError(s.T(), err)

	addr := s.redisServer.Addr()
	redisConn := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}

	s.accessRepository = NewRedisRepository(redisConn, s.log)
}

func (s *SuiteTestRepository) AfterTest(_, _ string) {
	s.SetupSuite()
	s.output = ""
}

func (s *SuiteTestRepository) TearDownSuite() {
	s.redisServer.Close()
}

func TestAccessRepository(t *testing.T) {
	suite.Run(t, new(SuiteTestRepository))
}

func (s *SuiteTestRepository) TestSet() {
	key := "don"
	value := "maria"
	ext := 8454
	err := s.accessRepository.Set(key, value, ext)
	require.NoError(s.T(), err)

	retValue, err := s.redisServer.Get(key)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), value, retValue)

	s.redisServer.FastForward(time.Second * 100000)

	_, err = s.redisServer.Get(key)
	assert.Equal(s.T(), err, miniredis.ErrKeyNotFound)
	assert.Equal(s.T(), s.output, "")

	s.redisServer.SetError("Error")
	err = s.accessRepository.Set(key, value, ext)
	assert.Error(s.T(), err)
	s.redisServer.Close()
}

func (s *SuiteTestRepository) TestGet() {
	key := "don"
	value := "maria"
	ext := 8454

	err := s.redisServer.Set(key, value)
	require.NoError(s.T(), err)
	s.redisServer.SetTTL(key, time.Duration(ext))

	var resVal string
	resVal, err = s.accessRepository.Get(key)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), resVal, value)
	assert.Equal(s.T(), s.output, "")

	s.redisServer.SetError("Error")
	_, err = s.accessRepository.Get(key)
	assert.Error(s.T(), err)

	s.redisServer.FastForward(time.Second * 100000)
	_, err = s.accessRepository.Get(key)
	assert.Error(s.T(), NotFound)

	s.redisServer.Close()
}

func (s *SuiteTestRepository) TestIncrment() {
	key := "don"
	value := "1"

	err := s.redisServer.Set(key, value)
	require.NoError(s.T(), err)

	var resValue int64
	resValue, err = s.accessRepository.Increment(key)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%d", resValue-1), value)

	s.redisServer.SetError("Error")
	_, err = s.accessRepository.Increment(key)
	assert.Error(s.T(), err)

	s.redisServer.Close()
}
