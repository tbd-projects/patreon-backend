package logout_handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type LogoutTestSuite struct {
	handlers.SuiteHandler
	handler *LogoutHandler
}

func (s *LogoutTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewLogoutHandler(s.Logger, s.MockSessionsManager)
}

func (s *LogoutTestSuite) TestPOST_WithSession() {
	uniqID := "1"
	test := handlers.TestTable{
		Name:              "with cookies",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}

	ctx := context.WithValue(context.Background(), "session_id", uniqID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(context.Background(), uniqID).
		Times(test.ExpectedMockTimes).
		Return(nil)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func (s *LogoutTestSuite) TestPOST_WithoutCookies() {
	uniqID := "1"
	test := handlers.TestTable{
		Name:              "without cookies",
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}

	reader, _ := http.NewRequest(http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(context.Background(), uniqID).
		Times(test.ExpectedMockTimes).
		Return(nil)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func (s *LogoutTestSuite) TestPOST_ErrorSessions() {
	uniqID := "1"
	test := handlers.TestTable{
		Name:              "without cookies",
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}

	ctx := context.WithValue(context.Background(), "session_id", uniqID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(context.Background(), uniqID).
		Times(test.ExpectedMockTimes).
		Return(errors.New(""))
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func TestLogoutSuite(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}
