package usecase_factory

import (
	"patreon/internal/app/delivery/http/handlers"
	mock_repository_factory "patreon/internal/app/usecase/usecase_factory/mocks"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type FactorySuite struct {
	handlers.SuiteHandler
	mockRepositoryFactory *mock_repository_factory.MockRepositoryFactory
	factory               *UsecaseFactory
}

func (s *FactorySuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.mockRepositoryFactory = mock_repository_factory.NewMockRepositoryFactory(s.Mock)
}
func (s *FactorySuite) TestGetUserUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetUserRepository()
	s.mockRepositoryFactory.EXPECT().GetFilesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetUserUsecase()
}
func (s *FactorySuite) TestGetUserUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.userUsecase = s.MockUserUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetUserUsecase()
}
func (s *FactorySuite) TestGetCreatorUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetCreatorRepository()
	s.mockRepositoryFactory.EXPECT().GetFilesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCreatorUsecase()
}
func (s *FactorySuite) TestGetCreatorUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.creatorUsecase = s.MockCreatorUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCreatorUsecase()
}
func (s *FactorySuite) TestGetCsrfrUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetCsrfRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCsrfUsecase()
}
func (s *FactorySuite) TestGetCsrfUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.csrfUsecase = s.MockCsrfUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetCsrfUsecase()
}
func (s *FactorySuite) TestGetAccessUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetAccessRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAccessUsecase()
}
func (s *FactorySuite) TestGetAccessUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.accessUsecase = s.MockAccessUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAccessUsecase()
}
func (s *FactorySuite) TestGetSubscribersUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
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
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.subscribersUsecase = s.MockSubscribersUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetSubscribersUsecase()
}

func (s *FactorySuite) TestGetAwardsUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.awardsUsecase = nil
	s.mockRepositoryFactory.EXPECT().GetAwardsRepository()
	s.mockRepositoryFactory.EXPECT().GetFilesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAwardsUsecase()
}

func (s *FactorySuite) TestGetAwardsUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.awardsUsecase = s.MockAwardsUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetAwardsUsecase()
}

func (s *FactorySuite) TestGetPostsUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetPostsRepository()
	s.mockRepositoryFactory.EXPECT().GetPostsDataRepository()
	s.mockRepositoryFactory.EXPECT().GetFilesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetPostsUsecase()
}

func (s *FactorySuite) TestGetPostsUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.postsUsecase = s.MockPostsUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetPostsUsecase()
}
func (s *FactorySuite) TestGetLikesUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetLikesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetLikesUsecase()
}

func (s *FactorySuite) TestGetLikesUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	factory.likesUsecase = s.MockLikeUsecase

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetLikesUsecase()
}
func (s *FactorySuite) TestGetPostsDataUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetPostsDataRepository()
	s.mockRepositoryFactory.EXPECT().GetFilesRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetPostsDataUsecase()
}

func (s *FactorySuite) TestGetInfoUsecaseFirstCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
	s.mockRepositoryFactory.EXPECT().GetInfoRepository()

	defer func() {
		if r := recover(); r != nil {
			assert.Fail(s.T(), "fail on getUserUsecase()")
		}
	}()
	factory.GetInfoUsecase()
}

func (s *FactorySuite) TestGetInfoUsecaseSecondCall() {
	factory := NewUsecaseFactory(s.mockRepositoryFactory)
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
