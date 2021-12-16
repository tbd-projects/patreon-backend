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
	mock_usecase "patreon/internal/app/usecase/awards/mocks"
	"strconv"
	"testing"
)

func TestAwardsMiddleware(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	mock := gomock.NewController(t)
	mockawards := mock_usecase.NewAwardsUsecase(mock)
	utilits := NewAwardsMiddleware(log, mockawards)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
		"award_id":   strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockawards.EXPECT().GetCreatorId(int64(1)).Return(int64(1), nil)
	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusOK)
	mock.Finish()
	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	vars = map[string]string{
		"creator_id": strconv.Itoa(1),
		"award_id":   strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockawards.EXPECT().GetCreatorId(int64(1)).Return(int64(2), nil)
	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusForbidden)
	mock.Finish()

	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	vars = map[string]string{
		"creator_id": strconv.Itoa(1),
		"award_id":   strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockawards.EXPECT().GetCreatorId(int64(1)).Return(int64(2), repository.DefaultErrDB)
	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	mock.Finish()

	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	vars = map[string]string{
		"creator_id": strconv.Itoa(1),
		"award_id":   strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	mockawards.EXPECT().GetCreatorId(int64(1)).Return(int64(2), repository.NotFound)
	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusForbidden)
	mock.Finish()

	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	vars = map[string]string{
		"award_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	mock.Finish()

	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	vars = map[string]string{
		"creator_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckCorrectAward(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)
	mock.Finish()
}
