package logout_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/models"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Data:              models.User{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.Data)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "session_id", uniqID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(uniqID).
		Times(test.ExpectedMockTimes).
		Return(nil)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func (s *LogoutTestSuite) TestPOST_WithoutCookies() {
	uniqID := "1"
	test := handlers.TestTable{
		Name:              "without cookies",
		Data:              &models.User{},
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.Data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(uniqID).
		Times(test.ExpectedMockTimes).
		Return(nil)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func (s *LogoutTestSuite) TestPOST_ErrorSessions() {
	uniqID := "1"
	test := handlers.TestTable{
		Name:              "without cookies",
		Data:              &models.User{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.Data)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "session_id", uniqID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/logout", &b)

	s.MockSessionsManager.EXPECT().
		Delete(uniqID).
		Times(test.ExpectedMockTimes).
		Return(errors.New(""))
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
}

func TestLogoutSuite(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}
