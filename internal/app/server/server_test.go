package server

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleRoot(t *testing.T) {
	expected := "hello patron!"

	server := New(NewConfig(), mux.NewRouter())
	recorder := httptest.NewRecorder()
	reader, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	server.HandleRoot().ServeHTTP(recorder, reader)

	assert.Equal(t, expected, recorder.Body.String())
}
