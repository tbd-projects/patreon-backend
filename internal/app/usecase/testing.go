package usecase

import (
	"io/ioutil"
	mock_repository_creator "patreon/internal/app/repository/creator/mocks"
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
	Mock                  *gomock.Controller
	MockCreatorRepository *mock_repository_creator.CreatorRepository
	MockUserRepository    *mock_repository_user.MockRepository

	Logger *logrus.Logger
	Tb     TestTable
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockCreatorRepository = mock_repository_creator.NewCreatorRepository(s.Mock)
	s.MockUserRepository = mock_repository_user.NewMockRepository(s.Mock)

	s.Logger = logrus.New()
	s.Logger.SetOutput(ioutil.Discard)

}
func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
