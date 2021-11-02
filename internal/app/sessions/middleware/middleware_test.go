package middleware

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/repository"
	mock_sessions "patreon/internal/app/sessions/mocks"
	"patreon/internal/app/sessions/models"
	"testing"
)

func TestSessionMiddleware_Check(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	sessionManager := mock_sessions.NewMockSessionsManager(mock)
	middleware := NewSessionMiddleware(sessionManager, log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	sessionId := "sadasd"
	cok := &http.Cookie{	}
	cok.Value = sessionId
	cok.Name = "session_id"
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.AddCookie(cok)
	require.NoError(t, err)
	res := models.Result{UserID: 1, UniqID: "asdasd"}
	sessionManager.EXPECT().Check(sessionId).Return(res, nil)
	middleware.Check(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdRaw := r.Context().Value("user_id")
		sessIdRaw := r.Context().Value("session_id")
		require.NotNil(t, userIdRaw)
		require.NotNil(t, sessIdRaw)
		userId, ok := userIdRaw.(int64)
		require.True(t, ok)
		sessId, ok := sessIdRaw.(string)
		require.True(t, ok)
		assert.Equal(t, userId, res.UserID)
		assert.Equal(t, sessId, res.UniqID)
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusOK)

	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	middleware.Check(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)

	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	reader.AddCookie(cok)
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	sessionManager.EXPECT().Check(sessionId).Return(res, repository.DefaultErrDB)
	middleware.Check(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestSessionMiddleware_CheckNotAuthorized(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	sessionManager := mock_sessions.NewMockSessionsManager(mock)
	middleware := NewSessionMiddleware(sessionManager, log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	sessionId := "sadasd"
	cok := &http.Cookie{	}
	cok.Value = sessionId
	cok.Name = "session_id"
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.AddCookie(cok)
	require.NoError(t, err)
	res := models.Result{UserID: 1, UniqID: "asdasd"}
	sessionManager.EXPECT().Check(sessionId).Return(res, nil)
	middleware.CheckNotAuthorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusTeapot)

	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	middleware.CheckNotAuthorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusOK)

	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	reader.AddCookie(cok)
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	sessionManager.EXPECT().Check(sessionId).Return(res, repository.DefaultErrDB)
	middleware.CheckNotAuthorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestSessionMiddleware_AddUserId(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	sessionManager := mock_sessions.NewMockSessionsManager(mock)
	middleware := NewSessionMiddleware(sessionManager, log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	sessionId := "sadasd"
	cok := &http.Cookie{	}
	cok.Value = sessionId
	cok.Name = "session_id"
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.AddCookie(cok)
	require.NoError(t, err)
	res := models.Result{UserID: 1, UniqID: "asdasd"}
	sessionManager.EXPECT().Check(sessionId).Return(res, nil)
	middleware.AddUserId(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdRaw := r.Context().Value("user_id")
		sessIdRaw := r.Context().Value("session_id")
		require.NotNil(t, userIdRaw)
		require.NotNil(t, sessIdRaw)
		userId, ok := userIdRaw.(int64)
		require.True(t, ok)
		sessId, ok := sessIdRaw.(string)
		require.True(t, ok)
		assert.Equal(t, userId, res.UserID)
		assert.Equal(t, sessId, res.UniqID)
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)
	assert.Equal(t, recorder.Code, http.StatusOK)

}