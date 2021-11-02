package usecase_csrf

import (
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/csrf/models"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	mock_repository_jwt "patreon/internal/app/csrf/repository/jwt/mocks"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/suite"
)

type SuiteCsrfUsecase struct {
	suite.Suite
	Mock               *gomock.Controller
	MockCsrfRepository *mock_repository_jwt.JwtRepository
	uc                 Usecase
}

func (s *SuiteCsrfUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockCsrfRepository = mock_repository_jwt.NewJwtRepository(s.Mock)
	s.uc = NewCsrfUsecase(s.MockCsrfRepository)
}

func TestSuiteCsrfUsecase(t *testing.T) {
	suite.Run(t, new(SuiteCsrfUsecase))
}

type SourcesWithMathcher struct {
	sources *models.TokenSources
}

func newSourcesWithMatcher(sources *models.TokenSources) gomock.Matcher {
	return &SourcesWithMathcher{sources: sources}

}
func (match *SourcesWithMathcher) String() string {
	return fmt.Sprintf("TokenSources: user_id: %v session_id: %v",
		match.sources.UserId, match.sources.SessionId)
}

func (match *SourcesWithMathcher) Matches(x interface{}) bool {
	switch x.(type) {
	case models.TokenSources:
		return x.(models.TokenSources).UserId == match.sources.UserId &&
			x.(models.TokenSources).SessionId == match.sources.SessionId
	default:
		return false
	}
}

func (s *SuiteCsrfUsecase) TestCsrfUsecase_Create_Ok() {
	sources := repository_jwt.TestSources(s.T())
	exp := models.Token("token")
	s.MockCsrfRepository.EXPECT().
		Create(newSourcesWithMatcher(sources)).
		Times(1).
		Return(exp, nil)

	token, err := s.uc.Create(sources.SessionId, sources.UserId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), exp, token)
}
func (s *SuiteCsrfUsecase) TestCsrfUsecase_Create_ErrorRepository() {
	sources := repository_jwt.TestSources(s.T())
	expErr := repository_jwt.ErrorSignedToken
	s.MockCsrfRepository.EXPECT().
		Create(newSourcesWithMatcher(sources)).
		Times(1).
		Return(models.Token(""), &app.GeneralError{Err: repository_jwt.ErrorSignedToken})

	token, err := s.uc.Create(sources.SessionId, sources.UserId)
	assert.Equal(s.T(), token, models.Token(""))
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expErr, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteCsrfUsecase) TestCsrfUsecase_Check_Ok() {
	sources := repository_jwt.TestSources(s.T())
	token := "token"
	s.MockCsrfRepository.EXPECT().
		Check(newSourcesWithMatcher(sources), models.Token(token)).
		Times(1).
		Return(nil)

	err := s.uc.Check(sources.SessionId, sources.UserId, token)
	assert.NoError(s.T(), err)
}
func (s *SuiteCsrfUsecase) TestCsrfUsecase_Check_BadToken() {
	sources := repository_jwt.TestSources(s.T())
	expErr := repository_jwt.BadToken
	token := "token"
	s.MockCsrfRepository.EXPECT().
		Check(newSourcesWithMatcher(sources), models.Token(token)).
		Times(1).
		Return(&app.GeneralError{Err: repository_jwt.BadToken})

	err := s.uc.Check(sources.SessionId, sources.UserId, token)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expErr, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteCsrfUsecase) TestCsrfUsecase_Check_TokenExpired() {
	sources := repository_jwt.TestSources(s.T())
	expErr := repository_jwt.TokenExpired
	token := "token"
	s.MockCsrfRepository.EXPECT().
		Check(newSourcesWithMatcher(sources), models.Token(token)).
		Times(1).
		Return(&app.GeneralError{Err: repository_jwt.TokenExpired})

	err := s.uc.Check(sources.SessionId, sources.UserId, token)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expErr, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteCsrfUsecase) TestCsrfUsecase_Check_ParseTokenError() {
	sources := repository_jwt.TestSources(s.T())
	expErr := repository_jwt.ParseClaimsError
	token := "token"
	s.MockCsrfRepository.EXPECT().
		Check(newSourcesWithMatcher(sources), models.Token(token)).
		Times(1).
		Return(&app.GeneralError{Err: repository_jwt.ParseClaimsError})

	err := s.uc.Check(sources.SessionId, sources.UserId, token)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expErr, errors.Cause(err).(*app.GeneralError).Err)
}
