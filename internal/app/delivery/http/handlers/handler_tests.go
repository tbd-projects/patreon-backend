package handlers

import (
	"github.com/gorilla/mux"
	"io"
	"patreon/internal/app"
	mock_sessions "patreon/internal/app/sessions/mocks"
	mock_usecase_awards "patreon/internal/app/usecase/awards/mocks"
	mock_usecase_creator "patreon/internal/app/usecase/creator/mocks"
	mock_usecase_user "patreon/internal/app/usecase/user/mocks"

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
	Mock                *gomock.Controller
	MockUserUsecase     *mock_usecase_user.MockUsecase
	MockCreatorUsecase  *mock_usecase_creator.CreatorUsecase
	MockAwardsUsecase   *mock_usecase_awards.AwardsUsecase
	MockSessionsManager *mock_sessions.MockSessionsManager
	Tb                  TestTable
	Logger              *logrus.Logger
	Router              *mux.Router
	Cors                *app.CorsConfig
}

func (s *SuiteHandler) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockUserUsecase = mock_usecase_user.NewMockUsecase(s.Mock)
	s.MockCreatorUsecase = mock_usecase_creator.NewCreatorUsecase(s.Mock)
	s.MockAwardsUsecase = mock_usecase_awards.NewAwardsUsecase(s.Mock)
	s.MockSessionsManager = mock_sessions.NewMockSessionsManager(s.Mock)

	s.Tb = TestTable{}
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteHandler) TearDownSuite() {
	s.Mock.Finish()
}
