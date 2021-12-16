package usecase_access

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/repository"
	repository_access "patreon/internal/app/repository/access"
	"patreon/internal/app/usecase"
	"testing"
)

type SuiteAccessUsecase struct {
	usecase.SuiteUsecase
	uc Usecase
}

func (s *SuiteAccessUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewAccessUsecase(s.MockAccessRepository)
}

func (s *SuiteAccessUsecase) TestCreatorUsecase_Create() {
	userIp := "123213"

	s.MockAccessRepository.EXPECT().
		Set(userIp, "0", timeLimit).
		Times(1).
		Return(nil)
	id, err := s.uc.Create(userIp)
	assert.True(s.T(), id)
	assert.NoError(s.T(), err)

	s.MockAccessRepository.EXPECT().
		Set(userIp, "0", timeLimit).
		Times(1).
		Return(repository.DefaultErrDB)
	id, err = s.uc.Create(userIp)
	assert.False(s.T(), id)
	assert.ErrorIs(s.T(), err, repository.DefaultErrDB)
}

func (s *SuiteAccessUsecase) TestCreatorUsecase_AddToBlackList() {
	userIp := "123213"

	s.MockAccessRepository.EXPECT().
		Set(blackList+userIp, userIp, timeBlocked).
		Times(1).
		Return(nil)
	err := s.uc.AddToBlackList(userIp)
	assert.NoError(s.T(), err)

	s.MockAccessRepository.EXPECT().
		Set(blackList+userIp, userIp, timeBlocked).
		Times(1).
		Return(repository.DefaultErrDB)
	err = s.uc.AddToBlackList(userIp)
	assert.ErrorIs(s.T(), err, repository.DefaultErrDB)
}

func (s *SuiteAccessUsecase) TestCreatorUsecase_CheckBlackList() {
	userIp := "123213"

	s.MockAccessRepository.EXPECT().
		Get(blackList+userIp).
		Times(1).
		Return("", nil)
	ok, err := s.uc.CheckBlackList(userIp)
	assert.True(s.T(), ok)
	assert.NoError(s.T(), err)

	s.MockAccessRepository.EXPECT().
		Get(blackList+userIp).
		Times(1).
		Return("", repository.DefaultErrDB)
	ok, err = s.uc.CheckBlackList(userIp)
	assert.True(s.T(), ok)
	assert.Error(s.T(), err)

	s.MockAccessRepository.EXPECT().
		Get(blackList+userIp).
		Times(1).
		Return("", repository_access.NotFound)
	ok, err = s.uc.CheckBlackList(userIp)
	assert.False(s.T(), ok)
	assert.NoError(s.T(), err)
}

func (s *SuiteAccessUsecase) TestCreatorUsecase_Update() {
	userIp := "123213"

	s.MockAccessRepository.EXPECT().
		Increment(userIp).
		Times(1).
		Return(int64(1), nil)
	id, err := s.uc.Update(userIp)
	assert.Equal(s.T(), int64(1), id)
	assert.NoError(s.T(), err)

	s.MockAccessRepository.EXPECT().
		Increment(userIp).
		Times(1).
		Return(int64(0), repository.DefaultErrDB)
	id, err = s.uc.Update(userIp)
	assert.Equal(s.T(), int64(-1), id)
	assert.ErrorIs(s.T(), err, repository.DefaultErrDB)

	s.MockAccessRepository.EXPECT().
		Increment(userIp).
		Times(1).
		Return(int64(queryLimit+10), nil)
	id, err = s.uc.Update(userIp)
	assert.Equal(s.T(), int64(-1), id)
	assert.ErrorIs(s.T(), err, NoAccess)
}

func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteAccessUsecase))
}
