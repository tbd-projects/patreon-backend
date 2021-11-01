package middleware

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUtilitiesMiddleware_CheckPanic(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	utilits := NewUtilitiesMiddleware(log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)

	utilits.CheckPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boba")
	})).ServeHTTP(recorder, reader)
}

func TestUtilitiesMiddleware_UpgradeLogger(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	utilits := NewUtilitiesMiddleware(log)

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)

	utilits.UpgradeLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := r.Context().Value("logger")
		require.NotNil(t, entry)
		entryParse, ok := entry.(*logrus.Entry)
		require.True(t, ok)
		assert.Equal(t, entryParse.Data["urls"], r.URL)
		assert.Equal(t, entryParse.Data["method"], r.Method)
		assert.Equal(t, entryParse.Data["remote_addr"], r.RemoteAddr)
		_, ok =  entryParse.Data["work_time"]
		assert.True(t,ok)
		_, ok =  entryParse.Data["req_id"]
		assert.True(t,ok)
	})).ServeHTTP(recorder, reader)
}