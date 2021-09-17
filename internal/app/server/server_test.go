package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleRoot(t *testing.T) {
	expected := "hello patron!"
	router := NewRouter()
	router.Configure()
	recorder := httptest.NewRecorder()
	reader, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	router.HandleRoot().ServeHTTP(recorder, reader)

	assert.Equal(t, expected, recorder.Body.String())
}
