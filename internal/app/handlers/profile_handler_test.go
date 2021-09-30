package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/store"
	"patreon/internal/models"
)

type ProfileTestSuite struct {
	SuiteTestBaseHandler
}


func (s *ProfileTestSuite) TestServeHTTP_Correct() {
	userID := int64(1)
	test := TestTable{
		name:              "correct",
		data:              &models.User{ID: int(userID), Login: "some", Nickname: "done"},
		expectedMockTimes: 1,
		expectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()
	handler := NewProfileHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)
	handler.SetStore(s.store)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", userID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/login", &b)

	s.mockUserRepository.EXPECT().FindByID(userID).Times(test.expectedMockTimes).Return(test.data.(*models.User), nil)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)

	req := &models.Profile{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &models.Profile{Nickname: test.data.(*models.User).Nickname,
		Avatar: test.data.(*models.User).Avatar})
}

func (s *ProfileTestSuite) TestServeHTTP_WitDBError() {
	userID := int64(1)
	test := TestTable{
		name:              "with db error",
		data:              nil,
		expectedMockTimes: 1,
		expectedCode:      http.StatusServiceUnavailable,
	}

	recorder := httptest.NewRecorder()
	handler := NewProfileHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)
	handler.SetStore(s.store)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", userID)
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/login", &b)

	s.mockUserRepository.EXPECT().FindByID(userID).Times(test.expectedMockTimes).Return(nil, store.NotFound)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)
}

func (s *ProfileTestSuite) TestServeHTTP_WithoutContext() {
	userID := int64(1)
	test := TestTable{
		name:              "without context",
		data:              nil,
		expectedMockTimes: 0,
		expectedCode:      http.StatusInternalServerError,
	}

	recorder := httptest.NewRecorder()
	handler := NewProfileHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)
	handler.SetStore(s.store)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)

	s.mockUserRepository.EXPECT().FindByID(userID).Times(test.expectedMockTimes)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)
}
