package usecase

import (
	"io/ioutil"
	"patreon/internal/app/models"
	mock_repository "patreon/internal/app/repository/creator/mocks"
	"testing"

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
	MockCreatorRepository *mock_repository.CreatorRepository
	Logger                *logrus.Logger
	Tb                    TestTable
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockCreatorRepository = mock_repository.NewCreatorRepository(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(ioutil.Discard)

}
func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}

func TestCreator(t *testing.T) *models.Creator {
	t.Helper()
	return &models.Creator{
		ID:          1,
		Category:    "podcasts",
		Nickname:    "podcaster2005",
		Description: "blog about IT",
	}
}
