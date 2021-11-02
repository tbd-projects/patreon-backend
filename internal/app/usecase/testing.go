package usecase

import (
	"io"
	mock_repository_awards "patreon/internal/app/repository/awards/mocks"
	mock_repository_creator "patreon/internal/app/repository/creator/mocks"
	mock_repository_files "patreon/internal/app/repository/files/mocks"
	mock_repository_likes "patreon/internal/app/repository/likes/mocks"
	mock_repository_posts "patreon/internal/app/repository/posts/mocks"
	mock_repository_posts_data "patreon/internal/app/repository/posts_data/mocks"
	mock_repository_subscribers "patreon/internal/app/repository/subscribers/mocks"
	mock_repository_user "patreon/internal/app/repository/user/mocks"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedError     error
}
type SuiteUsecase struct {
	suite.Suite
	Mock                      *gomock.Controller
	MockCreatorRepository     *mock_repository_creator.CreatorRepository
	MockUserRepository        *mock_repository_user.UserRepository
	MockSubscribersRepository *mock_repository_subscribers.SubscribersRepository
	MockFileRepository        *mock_repository_files.FilesRepository
	MockAwardsRepository      *mock_repository_awards.AwardsRepository
	MockPostsRepository       *mock_repository_posts.PostsRepository
	MockLikesRepository       *mock_repository_likes.LikesRepository
	MockPostsDataRepository   *mock_repository_posts_data.PostsDataRepository

	Logger *logrus.Logger
	Tb     TestTable
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockCreatorRepository = mock_repository_creator.NewCreatorRepository(s.Mock)
	s.MockUserRepository = mock_repository_user.NewUserRepository(s.Mock)
	s.MockSubscribersRepository = mock_repository_subscribers.NewSubscribersRepository(s.Mock)
	s.MockFileRepository = mock_repository_files.NewFilesRepository(s.Mock)
	s.MockPostsDataRepository = mock_repository_posts_data.NewPostsDataRepository(s.Mock)
	s.MockPostsRepository = mock_repository_posts.NewPostsRepository(s.Mock)
	s.MockAwardsRepository = mock_repository_awards.NewAwardsRepository(s.Mock)
	s.MockLikesRepository = mock_repository_likes.NewLikesRepository(s.Mock)

	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
