package usecase_subscribers

import (
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type SuiteSubscribersUsecase struct {
	usecase.SuiteUsecase
	uc Usecase
}

func (s *SuiteSubscribersUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewSubscribersUsecase(s.MockSubscribersRepository)
}

func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_OK() {
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber.UserID, subscriber.CreatorID).
		Return(false, nil).
		Times(1)
	s.MockSubscribersRepository.EXPECT().
		Create(subscriber).Times(1).
		Return(nil)
	err := s.uc.Subscribe(subscriber)
	assert.NoError(s.T(), err)
}
func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_AlreadyExists() {
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber.UserID, subscriber.CreatorID).
		Return(true, nil).
		Times(1)

	err := s.uc.Subscribe(subscriber)
	assert.Equal(s.T(), err, SubscriptionAlreadyExists)
}
func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_CheckExistsError() {
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber.UserID, subscriber.CreatorID).
		Return(false, repository.NewDBError(repository.DefaultErrDB)).
		Times(1)
	err := s.uc.Subscribe(subscriber)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteSubscribersUsecase) TestSubscribersUsecaseSubscribe_RepositoryCreateError() {
	subscriber := models.TestSubscriber()
	s.MockSubscribersRepository.EXPECT().
		Get(subscriber.UserID, subscriber.CreatorID).
		Return(false, nil).
		Times(1)
	s.MockSubscribersRepository.EXPECT().
		Create(subscriber).
		Times(1).
		Return(&app.GeneralError{
			Err: repository.DefaultErrDB,
		})
	err := s.uc.Subscribe(subscriber)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetCreators_OK() {
	subscriber := models.TestSubscriber()
	expCreators := []int64{1, 2}
	s.MockSubscribersRepository.EXPECT().
		GetCreators(subscriber.UserID).
		Times(1).
		Return(expCreators, nil)
	res, err := s.uc.GetCreators(subscriber.UserID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expCreators, res)
}
func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetCreators_RepositoryError() {
	subscriber := models.TestSubscriber()
	expCreators := []int64{}
	s.MockSubscribersRepository.EXPECT().
		GetCreators(subscriber.UserID).
		Times(1).
		Return(expCreators, repository.NewDBError(repository.DefaultErrDB))
	res, err := s.uc.GetCreators(subscriber.UserID)
	assert.Equal(s.T(), expCreators, res)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetSubscribers_OK() {
	subscriber := models.TestSubscriber()
	expUsers := []int64{1, 2}
	s.MockSubscribersRepository.EXPECT().
		GetSubscribers(subscriber.CreatorID).
		Times(1).
		Return(expUsers, nil)
	res, err := s.uc.GetSubscribers(subscriber.CreatorID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expUsers, res)
}
func (s *SuiteSubscribersUsecase) TestSubscriberUsecaseGetSubscribers_RepositoryError() {
	subscriber := models.TestSubscriber()
	expUsers := []int64{}
	s.MockSubscribersRepository.EXPECT().
		GetSubscribers(subscriber.CreatorID).
		Times(1).
		Return(expUsers, repository.NewDBError(repository.DefaultErrDB))

	res, err := s.uc.GetSubscribers(subscriber.CreatorID)
	assert.Equal(s.T(), expUsers, res)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.DefaultErrDB, errors.Cause(err).(*app.GeneralError).Err)
}
func TestSubscribersUsecase(t *testing.T) {
	suite.Run(t, new(SuiteSubscribersUsecase))
}
