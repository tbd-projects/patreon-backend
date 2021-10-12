package usecase_creator

import (
	"io/ioutil"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	mock_repository "patreon/internal/app/repository/creator/mocks"
	mock_sessions "patreon/internal/app/sessions/mocks"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

type TestTable struct {
	name              string
	data              interface{}
	expectedMockTimes int
	expectedError     error
}
type SuiteCreatorUsecase struct {
	suite.Suite
	mock                  *gomock.Controller
	mockCreatorRepository *mock_repository.CreatorRepository
	mockSessionManager    *mock_sessions.MockSessionsManager
	uc                    Usecase
	logger                *logrus.Logger
}

func (s *SuiteCreatorUsecase) SetupSuite() {
	s.mock = gomock.NewController(s.T())
	s.mockCreatorRepository = mock_repository.NewCreatorRepository(s.mock)
	s.mockSessionManager = mock_sessions.NewMockSessionsManager(s.mock)
	s.logger = logrus.New()
	s.logger.SetOutput(ioutil.Discard)
}
func (s *SuiteCreatorUsecase) TearDownSuite() {
	s.mock.Finish()
}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_DB_Error() {
	test := TestTable{
		name:              "Repository return error",
		data:              TestCreator(s.T()),
		expectedMockTimes: 1,
		expectedError:     repository.DefaultErrDB,
	}
	uc := NewCreatorUsecase(s.mockCreatorRepository)
	cr := TestCreator(s.T())
	s.mockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(test.expectedMockTimes).
		Return(cr, repository.DefaultErrDB)

	id, err := uc.Create(cr)

	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), test.expectedError, errors.Cause(err))
}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Creator_Exist() {
	test := TestTable{
		name:              "Creator already exist",
		data:              TestCreator(s.T()),
		expectedMockTimes: 1,
		expectedError:     CreatorExist,
	}
	uc := NewCreatorUsecase(s.mockCreatorRepository)
	cr := TestCreator(s.T())
	s.mockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(test.expectedMockTimes).
		Return(cr, nil)
	id, err := uc.Create(cr)
	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), test.expectedError, err)

}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Validate_Error() {
	test := TestTable{
		name:              "Invalid creator data",
		data:              TestCreator(s.T()),
		expectedMockTimes: 1,
		expectedError:     models.IncorrectCreatorNickname,
	}
	uc := NewCreatorUsecase(s.mockCreatorRepository)
	cr := TestCreator(s.T())
	cr.Nickname = ""
	s.mockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(test.expectedMockTimes).
		Return(nil, nil)
	id, err := uc.Create(cr)
	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), test.expectedError, err)
}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Success() {
	test := TestTable{
		name:              "Success create creator",
		data:              TestCreator(s.T()),
		expectedMockTimes: 1,
		expectedError:     nil,
	}
	uc := NewCreatorUsecase(s.mockCreatorRepository)
	cr := TestCreator(s.T())
	s.mockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(test.expectedMockTimes).
		Return(nil, nil)

	s.mockCreatorRepository.EXPECT().
		Create(cr).
		Times(test.expectedMockTimes).
		Return(cr.ID, nil)

	id, err := uc.Create(cr)
	expectId := int64(1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), test.expectedError, err)
}

func TestUsecases(t *testing.T) {
	suite.Run(t, new(SuiteCreatorUsecase))
}
