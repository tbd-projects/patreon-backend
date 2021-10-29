package repository_jwt

import (
	"patreon/internal/app"
	"patreon/internal/app/csrf/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type SuiteJwtRepository struct {
	suite.Suite
	repository *JwtRepository
}

func (s *SuiteJwtRepository) SetupSuite() {
	s.repository = NewJwtRepository()
}
func (s *SuiteJwtRepository) TearDownSuite() {
}

func TestJwtRepository(t *testing.T) {
	suite.Run(t, new(SuiteJwtRepository))
}

func (s *SuiteJwtRepository) TestJwtRepository_Create_Ok() {
	sourses := *TestSources(s.T())
	res, err := s.repository.Create(sourses)
	assert.NoError(s.T(), err)
	err = s.repository.Check(sourses, res)
	assert.NoError(s.T(), err)
}
func (s *SuiteJwtRepository) TestJwtRepository_Check_TokenExpire() {
	sourses := *TestSources(s.T())
	expectedError := TokenExpired
	sourses.ExpiredTime = time.Now().Add(-1 * time.Hour)
	res, err := s.repository.Create(sourses)
	assert.NoError(s.T(), err)
	err = s.repository.Check(sourses, res)
	var generalError *app.GeneralError

	assert.True(s.T(), errors.As(err, &generalError))
	assert.Equal(s.T(), errors.Cause(err).(*app.GeneralError).Err, expectedError)
}
func (s *SuiteJwtRepository) TestJwtRepository_Check_IncorrectEncoding() {
	sourses := *TestSources(s.T())
	jwtClaimsTest := jwtCsrfClaims{
		SessionId: sourses.SessionId,
		UserId:    sourses.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: sourses.ExpiredTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwtClaimsTest)
	testToken, err := token.SignedString(s.repository.Secret)
	assert.NoError(s.T(), err)

	err = s.repository.Check(sourses, models.Token(testToken))

	var generalError *app.GeneralError
	assert.True(s.T(), errors.As(err, &generalError))
	assert.Equal(s.T(), errors.Cause(err).(*app.GeneralError).ExternalErr, IncorrectTokenSigningMethod)
}
