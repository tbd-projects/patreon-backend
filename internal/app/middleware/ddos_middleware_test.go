package middleware

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app"
	"patreon/internal/app/repository"
	usecase_access "patreon/internal/app/usecase/access"
	mock_usecase "patreon/internal/app/usecase/access/mocks"
	"testing"
)

type SuiteDdosMiddleware struct {
	suite.Suite
	mock       *gomock.Controller
	ip         string
	mockaccess *mock_usecase.AccessUsecase
	ddos DDosMiddleware
}

func (s *SuiteDdosMiddleware) SetupSuite() {
	log := &logrus.Logger{}
	s.mock = gomock.NewController(s.T())
	s.mockaccess = mock_usecase.NewAccessUsecase(s.mock)
	s.ddos = NewDdosMiddleware(log, s.mockaccess)
	s.ip = "123123"
}

func (s* SuiteDdosMiddleware) fterTest(_, _ string) {
	s.mock.Finish()
}


func (s *SuiteDdosMiddleware) TestPostsMiddleware_OK() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(true, nil)
	s.mockaccess.EXPECT().Update(s.ip).Return(int64(1), nil)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusOK)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_NoAccess() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(false, usecase_access.NoAccess)
	s.mockaccess.EXPECT().AddToBlackList(s.ip).Return(nil)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusTooManyRequests)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_NoAccessError() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(false, repository.DefaultErrDB)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusInternalServerError)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_InBlackList() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(true, nil)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusTooManyRequests)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_FirstQuery() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(true, usecase_access.FirstQuery)
	s.mockaccess.EXPECT().Create(s.ip).Return(true, nil)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusOK)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_FirstQuery_Error() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(true, usecase_access.FirstQuery)
	s.mockaccess.EXPECT().Create(s.ip).Return(false, repository.DefaultErrDB)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusInternalServerError)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_NoAccess2() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(true, nil)
	s.mockaccess.EXPECT().Update(s.ip).Return(int64(app.InvalidInt), usecase_access.NoAccess)
	s.mockaccess.EXPECT().AddToBlackList(s.ip).Return(nil)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusTooManyRequests)
}

func (s *SuiteDdosMiddleware) TestPostsMiddleware_ErrorBd() {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(s.T())


	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	reader.RemoteAddr = s.ip
	require.NoError(s.T(), err)

	s.mockaccess.EXPECT().CheckBlackList(s.ip).Return(false, nil)
	s.mockaccess.EXPECT().CheckAccess(s.ip).Return(true, nil)
	s.mockaccess.EXPECT().Update(s.ip).Return(int64(app.InvalidInt), repository.DefaultErrDB)
	s.ddos.CheckAccess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(s.T(), recorder.Code, http.StatusInternalServerError)
}


func TestDdosMiddleware(t *testing.T) {
	suite.Run(t, new(SuiteDdosMiddleware))
}
