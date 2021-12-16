package csrf_handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app"
	"patreon/internal/app/csrf/csrf_models"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	"patreon/internal/app/delivery/http/handlers"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type CsrfTestSuite struct {
	handlers.SuiteHandler
	handler *CsrfHandler
}

func (s *CsrfTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewCsrfHandler(s.Logger, s.MockSessionsManager, s.MockCsrfUsecase)
}
func TestCsrfHandler(t *testing.T) {
	suite.Run(t, new(CsrfTestSuite))
}
func (s *CsrfTestSuite) TestCsrfHandlerGet_ServerErrorSession() {
	s.Tb = handlers.TestTable{
		Name:              "server error - invalid session_id",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	w := httptest.NewRecorder()

	b := bytes.Buffer{}

	r, _ := http.NewRequest(http.MethodGet, "/token", &b)
	r = r.WithContext(context.WithValue(r.Context(), "invalid session_id", "empty"))
	s.handler.GET(w, r)
	assert.Equal(s.T(), s.Tb.ExpectedCode, w.Code)
}
func (s *CsrfTestSuite) TestCsrfHandlerGet_ServerErrorUserId() {
	s.Tb = handlers.TestTable{
		Name:              "server error - invalid user_id",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	w := httptest.NewRecorder()

	b := bytes.Buffer{}

	r, _ := http.NewRequest(http.MethodGet, "/token", &b)
	r = r.WithContext(context.WithValue(r.Context(), "session_id", "session"))
	r = r.WithContext(context.WithValue(r.Context(), "invalid_user_id", "empty"))

	s.handler.GET(w, r)
	assert.Equal(s.T(), s.Tb.ExpectedCode, w.Code)
}
func (s *CsrfTestSuite) TestCsrfHandlerGet_CreateTokenError() {
	s.Tb = handlers.TestTable{
		Name:              "server error - error on create token",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	w := httptest.NewRecorder()

	b := bytes.Buffer{}

	r, _ := http.NewRequest(http.MethodGet, "/token", &b)
	session_id := "session"
	user_id := int64(1)
	r = r.WithContext(context.WithValue(r.Context(), "session_id", session_id))
	r = r.WithContext(context.WithValue(r.Context(), "user_id", user_id))

	s.MockCsrfUsecase.EXPECT().
		Create(session_id, user_id).
		Times(1).
		Return(csrf_models.Token(""), &app.GeneralError{
			Err: repository_jwt.ErrorSignedToken,
		})

	s.handler.GET(w, r)
	assert.Equal(s.T(), s.Tb.ExpectedCode, w.Code)
}
func (s *CsrfTestSuite) TestCsrfHandlerGet_Ok() {
	s.Tb = handlers.TestTable{
		Name:              "OK",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	w := httptest.NewRecorder()

	b := bytes.Buffer{}

	r, _ := http.NewRequest(http.MethodGet, "/token", &b)
	session_id := "session"
	user_id := int64(1)
	r = r.WithContext(context.WithValue(r.Context(), "session_id", session_id))
	r = r.WithContext(context.WithValue(r.Context(), "user_id", user_id))

	s.MockCsrfUsecase.EXPECT().
		Create(session_id, user_id).
		Times(1).
		Return(csrf_models.Token("token"), nil)

	s.handler.GET(w, r)
	assert.Equal(s.T(), s.Tb.ExpectedCode, w.Code)
}
