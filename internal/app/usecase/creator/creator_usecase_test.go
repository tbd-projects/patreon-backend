package usecase_creator

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"io"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	repository_files "patreon/internal/microservices/files/files/repository/files"
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
	s.uc = NewCreatorUsecase(s.MockCreatorRepository, s.MockFileClient, s.MockConvector)
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
		ExistsCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(false, repository.DefaultErrDB)
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
		ExistsCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(true, nil)
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
		ExistsCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(false, repository.NotFound)
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
		ExistsCreator(cr.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(false, repository.NotFound)

	s.MockCreatorRepository.EXPECT().
		Create(cr).
		Times(s.Tb.ExpectedMockTimes).
		Return(cr.ID, nil)

	id, err := s.uc.Create(cr)
	assert.Equal(s.T(), expectId, id)
	assert.Equal(s.T(), s.Tb.ExpectedError, err)
}

func (s *SuiteCreatorUsecase) TestCreatorUsecase_UpdateAvatar_Success() {
	cr := models.TestCreator()
	name := "true"
	out := io.Reader(bytes.NewBufferString(""))

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), out, repository_files.FileName(name)).
		Times(1).
		Return(out, repository_files.FileName(name), nil)

	s.MockCreatorRepository.EXPECT().
		ExistsCreator(cr.ID).
		Times(1).
		Return(true, nil)

	s.MockFileClient.EXPECT().
		SaveFile(context.Background(), out, repository_files.FileName(name), repository_files.Image).
		Times(1).
		Return(name, nil)

	s.MockCreatorRepository.EXPECT().
		UpdateAvatar(cr.ID, app.LoadFileUrl+name).
		Times(1).
		Return(nil)

	err := s.uc.UpdateAvatar(out, repository_files.FileName(name), cr.ID)
	assert.NoError(s.T(), err)
}

func (s *SuiteCreatorUsecase) TestCreatorUsecase_UpdateCover_Success() {
	cr := models.TestCreator()
	name := "true"
	out := io.Reader(bytes.NewBufferString(""))

	s.MockConvector.EXPECT().
		Convert(gomock.Any(), out, repository_files.FileName(name)).
		Times(1).
		Return(out, repository_files.FileName(name), nil)

	s.MockCreatorRepository.EXPECT().
		ExistsCreator(cr.ID).
		Times(1).
		Return(true, nil)

	s.MockFileClient.EXPECT().
		SaveFile(context.Background(), out, repository_files.FileName(name), repository_files.Image).
		Times(1).
		Return(name, nil)

	s.MockCreatorRepository.EXPECT().
		UpdateCover(cr.ID, app.LoadFileUrl+name).
		Times(1).
		Return(nil)

	err := s.uc.UpdateCover(out, repository_files.FileName(name), cr.ID)
	assert.NoError(s.T(), err)
}

func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteCreatorUsecase))
}
