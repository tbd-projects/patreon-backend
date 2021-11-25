package usecase_info

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/usecase"
	"testing"
)

type SuiteLikeUsecase struct {
	usecase.SuiteUsecase
	uc    Usecase
	testData *models.Info
}

func (s *SuiteLikeUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewInfoUsecase(s.MockInfoRepository)
	s.testData = &models.Info{TypePostData: []string{"don", "con"}, Category: []string{"gfy", "ton"}}
}

func (s *SuiteLikeUsecase) TestInfoUsecase_Get_Ok() {
	s.MockInfoRepository.EXPECT().
		Get().
		Times(1).
		Return(s.testData, nil)
	res, err := s.uc.Get()
	assert.Equal(s.T(), res, s.testData)
	assert.NoError(s.T(), err)
}

func (s *SuiteLikeUsecase) TestInfoUsecase_Get_Error() {
	s.MockInfoRepository.EXPECT().
		Get().
		Times(1).
		Return(nil, repository.DefaultErrDB)
	_, err := s.uc.Get()
	assert.Error(s.T(), err, repository.DefaultErrDB)
}

func TestUsecaseInfo(t *testing.T) {
	suite.Run(t, new(SuiteLikeUsecase))
}



