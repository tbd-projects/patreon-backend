package handler_factory

import (
	mock_usecase_factory "patreon/internal/app/delivery/http/handler_factory/mocks"
	"patreon/internal/app/delivery/http/handlers"
	mock_auth_checker "patreon/internal/microservices/auth/delivery/grpc/client/mocks"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type FactorySuite struct {
	handlers.SuiteHandler
	usecaseFactory *mock_usecase_factory.MockUsecaseFactory
	sessionService *mock_auth_checker.MockAuthCheckerClient
	factory        *HandlerFactory
}

func (s *FactorySuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.usecaseFactory = mock_usecase_factory.NewMockUsecaseFactory(s.Mock)
	s.sessionService = mock_auth_checker.NewMockAuthCheckerClient(s.Mock)
	sessionCon := &grpc.ClientConn{}
	s.factory = NewFactory(s.Logger, s.usecaseFactory, sessionCon)
}

func (s *FactorySuite) TestInitHandlers() {
	s.usecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetAttachesUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPaymentsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetInfoUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetStatsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCommentsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPayTokenUsecase().Times(1)

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
	s.usecaseFactory.EXPECT().GetUserUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCreatorUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCsrfUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetAwardsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPostsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetLikesUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetSubscribersUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetAttachesUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPaymentsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetInfoUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetStatsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetCommentsUsecase().Times(1)
	s.usecaseFactory.EXPECT().GetPayTokenUsecase().Times(1)

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
