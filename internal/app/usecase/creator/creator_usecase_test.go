package usecase_creator

import (
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type SuiteCreatorUsecase struct {
	usecase.SuiteUsecase
	uc Usecase
}

func (s *SuiteCreatorUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewCreatorUsecase(s.MockCreatorRepository, s.MockFilesRepository)
}

func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_DB_Error() {
	s.Tb = usecase.TestTable{
		Name:              "Repository return error",
		Data:              models.TestCreator(),
		ExpectedMockTimes: 1,
		ExpectedError:     repository.DefaultErrDB,
	}
	cr := models.TestCreator()
	s.MockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(cr, repository.DefaultErrDB)

	id, err := s.uc.Create(cr)

	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), s.Tb.ExpectedError, errors.Cause(err))
}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Creator_Exist() {
	s.Tb = usecase.TestTable{
		Name:              "Creator already exist",
		Data:              models.TestCreator(),
		ExpectedMockTimes: 1,
		ExpectedError:     CreatorExist,
	}
	cr := models.TestCreator()
	s.MockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(cr, nil)
	id, err := s.uc.Create(cr)
	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)

}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Validate_Error() {
	s.Tb = usecase.TestTable{
		Name:              "Invalid creator data",
		Data:              models.TestCreator(),
		ExpectedMockTimes: 1,
		ExpectedError:     models.IncorrectCreatorNickname,
	}
	cr := models.TestCreator()
	cr.Nickname = ""
	s.MockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	id, err := s.uc.Create(cr)
	expectId := int64(-1)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)
}
func (s *SuiteCreatorUsecase) TestCreatorUsecase_Create_Success() {
	s.Tb = usecase.TestTable{
		Name:              "Success create creator",
		Data:              models.TestCreator(),
		ExpectedMockTimes: 1,
		ExpectedError:     nil,
	}
	cr := models.TestCreator()
	expectId := cr.ID

	s.MockCreatorRepository.EXPECT().
		GetCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)

	s.MockCreatorRepository.EXPECT().
		Create(cr).
		Times(s.Tb.ExpectedMockTimes).
		Return(cr.ID, nil)

	id, err := s.uc.Create(cr)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)
}

func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteCreatorUsecase))
}
