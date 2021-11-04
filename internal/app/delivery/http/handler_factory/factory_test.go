package handler_factory

import (
	mock_usecase_factory "patreon/internal/app/delivery/http/handler_factory/mocks"
	"patreon/internal/app/delivery/http/handlers"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type FactorySuite struct {
	handlers.SuiteHandler
	CsrfUsecaseFactory *mock_usecase_factory.MockUsecaseFactory
	factory            *HandlerFactory
}

func (s *FactorySuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.CsrfUsecaseFactory = mock_usecase_factory.NewMockUsecaseFactory(s.Mock)
	s.factory = NewFactory(s.Logger, s.CsrfUsecaseFactory)
}

func (s *FactorySuite) TestInitHandlers() {
	s.CsrfUsecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetSessionManager().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetPostsDataUsecase().Times(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on initAllHandlers")
		}
	}()
	s.factory.initAllHandlers()
}
func (s *FactorySuite) TestGetHandlersUrlsFirstRun() {
	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on initAllHandlers")
		}
	}()
	s.factory.GetHandleUrls()
}
func (s *FactorySuite) TestGetHandlersUrlsAlreadyExists() {
	s.CsrfUsecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetSessionManager().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.CsrfUsecaseFactory.EXPECT().GetPostsDataUsecase().Times(1)
	s.factory.urlHandler = nil
	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on initAllHandlers")
		}
	}()
	s.factory.GetHandleUrls()
}
func TestFactoryHandler(t *testing.T) {
	suite.Run(t, new(FactorySuite))
}
