package handler_factory

import (
	"patreon/internal/app"
	mock_usecase_factory "patreon/internal/app/delivery/http/handler_factory/mocks"
	"patreon/internal/app/delivery/http/handlers"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type FactorySuite struct {
	handlers.SuiteHandler
	mockUsecaseFactory *mock_usecase_factory.MockUsecaseFactory
	factory            *HandlerFactory
}

func (s *FactorySuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.mockUsecaseFactory = mock_usecase_factory.NewMockUsecaseFactory(s.Mock)
	s.factory = NewFactory(s.Logger, s.Router, s.Cors, s.mockUsecaseFactory)
}

func (s *FactorySuite) TestInitHandlers() {
	s.mockUsecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetSessionManager().Times(1)
	s.mockUsecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetPostsDataUsecase().Times(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on initAllHandlers")
		}
	}()
	s.factory.initAllHandlers()
}
func (s *FactorySuite) TestGetHandlersUrlsFirstRun() {
	s.mockUsecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetSessionManager().Times(1)
	s.mockUsecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.mockUsecaseFactory.EXPECT().GetPostsDataUsecase().Times(1)

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on initAllHandlers")
		}
	}()
	s.factory.GetHandleUrls()
}
func (s *FactorySuite) TestGetHandlersUrlsAlreadyExists() {

	s.factory.urlHandler = &map[string]app.Handler{}
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
