package middleware

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreatorsMiddleware(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	utilits := NewCreatorsMiddleware(log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), "user_id", int64(1))
	reader, err := http.NewRequestWithContext(ctx, http.MethodPost, "/register", &b)
	vars := map[string]string{
		"creator_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckAllowUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusOK)

	recorder = httptest.NewRecorder()
	ctx = context.WithValue(context.Background(), "user_id", int64(2))
	reader, err = http.NewRequestWithContext(ctx, http.MethodPost, "/register", &b)
	vars = map[string]string{
		"creator_id": strconv.Itoa(1),
	}
	reader = mux.SetURLVars(reader, vars)
	require.NoError(t, err)

	utilits.CheckAllowUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusForbidden)

	recorder = httptest.NewRecorder()
	ctx = context.WithValue(context.Background(), "user_id", int64(2))
	reader, err = http.NewRequestWithContext(ctx, http.MethodPost, "/register", &b)
	require.NoError(t, err)

	utilits.CheckAllowUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusBadRequest)

	recorder = httptest.NewRecorder()
	reader, err = http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)

	utilits.CheckAllowUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, reader)

	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}
