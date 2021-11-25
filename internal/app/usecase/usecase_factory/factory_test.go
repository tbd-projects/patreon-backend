package usecase_factory

import (
	"patreon/internal/app/delivery/http/handlers"
	mock_repository_factory "patreon/internal/app/usecase/usecase_factory/mocks"
	"testing"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	handlers.SuiteHandler
	mockRepositoryFactory *mock_repository_factory.MockRepositoryFactory
	factory               *UsecaseFactory
	fileConn              *grpc.ClientConn
}

func (s *FactorySuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.mockRepositoryFactory = mock_repository_factory.NewMockRepositoryFactory(s.Mock)
	s.fileConn, _ = grpc.Dial("", grpc.WithInsecure())
}
func (s *FactorySuite) TestGetUserUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetUserRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetUserUsecase()
}
func (s *FactorySuite) TestGetUserUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.userUsecase = s.MockUserUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetUserUsecase()
}
func (s *FactorySuite) TestGetCreatorUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetCreatorRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCreatorUsecase()
}
func (s *FactorySuite) TestGetCreatorUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.creatorUsecase = s.MockCreatorUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCreatorUsecase()
}
func (s *FactorySuite) TestGetCsrfrUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetCsrfRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCsrfUsecase()
}
func (s *FactorySuite) TestGetCsrfUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.csrfUsecase = s.MockCsrfUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCsrfUsecase()
}
func (s *FactorySuite) TestGetAccessUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	s.mockRepositoryFactory.EXPECT().GetAccessRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAccessUsecase()
}
func (s *FactorySuite) TestGetAccessUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	factory.accessUsecase = s.MockAccessUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAccessUsecase()
}
func (s *FactorySuite) TestGetSubscribersUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	s.mockRepositoryFactory.EXPECT().GetSubscribersRepository()
	s.mockRepositoryFactory.EXPECT().GetAwardsRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetSubscribersUsecase()
}

func (s *FactorySuite) TestGetSubscribersUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	factory.subscribersUsecase = s.MockSubscribersUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetSubscribersUsecase()
}

func (s *FactorySuite) TestGetAwardsUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	factory.awardsUsecase = nil
	s.mockRepositoryFactory.EXPECT().GetAwardsRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAwardsUsecase()
}

func (s *FactorySuite) TestGetAwardsUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.awardsUsecase = s.MockAwardsUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAwardsUsecase()
}

func (s *FactorySuite) TestGetPostsUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)

	s.mockRepositoryFactory.EXPECT().GetPostsRepository()
	s.mockRepositoryFactory.EXPECT().GetAttachesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetPostsUsecase()
}

func (s *FactorySuite) TestGetPostsUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.postsUsecase = s.MockPostsUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetPostsUsecase()
}
func (s *FactorySuite) TestGetLikesUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetLikesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetLikesUsecase()
}

func (s *FactorySuite) TestGetLikesUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.likesUsecase = s.MockLikeUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetLikesUsecase()
}
func (s *FactorySuite) TestGetAttachesUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetAttachesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAttachesUsecase()
}

func (s *FactorySuite) TestGetInfoUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	s.mockRepositoryFactory.EXPECT().GetInfoRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetInfoUsecase()
}

func (s *FactorySuite) TestGetInfoUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory, s.fileConn)
	factory.infoUsecase = s.MockInfoUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetInfoUsecase()
}

func TestFactoryHandler(t *testing.T) {
	suite.Run(t, new(FactorySuite))
}
