package middleware

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/repository"
	mock_usecase "patreon/internal/app/usecase/posts/mocks"
	"strconv"
	"testing"
)

func TestPostsMiddleware(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
		"post_id":    strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockposts.EXPECT().GetCreatorId(int64(1)).Return(int64(1), nil)
	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusOK)
	mock.Finish()
}

func TestPostsMiddleware_Fobbiden(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)

	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
		"post_id":    strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockposts.EXPECT().GetCreatorId(int64(1)).Return(int64(2), nil)
	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusForbidden)
	mock.Finish()
}

func TestPostsMiddleware_StatusInternalServerError(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
		"post_id":    strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockposts.EXPECT().GetCreatorId(int64(1)).Return(int64(2), repository.DefaultErrDB)
	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	mock.Finish()
}

func TestPostsMiddleware_StatusForbidden2(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
		"post_id":    strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockposts.EXPECT().GetCreatorId(int64(1)).Return(int64(2), repository.NotFound)
	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusForbidden)
	mock.Finish()
}

func TestPostsMiddleware_StatusBadRequest(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"post_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	mock.Finish()
}

func TestPostsMiddleware_StatusBadRequest2(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockposts := mock_usecase.NewPostsUsecase(mock)
	utilits := NewPostsMiddleware(log, mockposts)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckCorrectPost(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	mock.Finish()
}
