package handlers

import (
	"io"
	mock_usecase_csrf "patreon/internal/app/csrf/usecase/mocks"
	mock_usecase "patreon/internal/app/usecase/access/mocks"
	mock_usecase_awards "patreon/internal/app/usecase/awards/mocks"
	mock_usecase_creator "patreon/internal/app/usecase/creator/mocks"
	mock_usecase_info "patreon/internal/app/usecase/info/mocks"
	mock_usecase_like "patreon/internal/app/usecase/likes/mocks"
	mock_usecase_posts "patreon/internal/app/usecase/posts/mocks"
	mock_usecase_attaches "patreon/internal/app/usecase/attaches/mocks"
	mock_subscribers "patreon/internal/app/usecase/subscribers/mocks"
	mock_usecase_user "patreon/internal/app/usecase/user/mocks"
	mock_auth_checker "patreon/internal/microservices/auth/delivery/grpc/client/mocks"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedCode      int
}

type SuiteHandler struct {
	suite.Suite
	Mock                   *gomock.Controller
	MockUserUsecase        *mock_usecase_user.UserUsecase
	MockLikeUsecase        *mock_usecase_like.LikesUsecase
	MockCreatorUsecase     *mock_usecase_creator.CreatorUsecase
	MockAwardsUsecase      *mock_usecase_awards.AwardsUsecase
	MockPostsUsecase       *mock_usecase_posts.PostsUsecase
	MockAttachesUsecase   *mock_usecase_attaches.AttachesUsecase
	MockAccessUsecase      *mock_usecase.AccessUsecase
	MockSessionsManager    *mock_auth_checker.MockAuthCheckerClient
	MockInfoUsecase        *mock_usecase_info.InfoUsecase
	Tb                     TestTable
	Logger                 *logrus.Logger
	MockCsrfUsecase        *mock_usecase_csrf.CsrfUsecase
	MockSubscribersUsecase *mock_subscribers.SubscribersUsecase
}

func (s *SuiteHandler) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockUserUsecase = mock_usecase_user.NewUserUsecase(s.Mock)
	s.MockCreatorUsecase = mock_usecase_creator.NewCreatorUsecase(s.Mock)
	s.MockAwardsUsecase = mock_usecase_awards.NewAwardsUsecase(s.Mock)
	s.MockSessionsManager = mock_auth_checker.NewMockAuthCheckerClient(s.Mock)
	s.MockCsrfUsecase = mock_usecase_csrf.NewCsrfUsecase(s.Mock)
	s.MockSubscribersUsecase = mock_subscribers.NewSubscribersUsecase(s.Mock)
	s.MockAccessUsecase = mock_usecase.NewAccessUsecase(s.Mock)
	s.MockLikeUsecase = mock_usecase_like.NewLikesUsecase(s.Mock)
	s.MockPostsUsecase = mock_usecase_posts.NewPostsUsecase(s.Mock)
	s.MockAttachesUsecase = mock_usecase_attaches.NewAttachesUsecase(s.Mock)
	s.MockInfoUsecase = mock_usecase_info.NewInfoUsecase(s.Mock)

	s.Tb = TestTable{}
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteHandler) TearDownSuite() {
	s.Mock.Finish()
}
