package utilits

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/models"
	"testing"
)

func TestResponder(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	responder := Responder{NewLogObject(log)}

	b := bytes.Buffer{}
	recorder := httptest.NewRecorder()
	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)

	tmpError := errors.New("some error")
	responder.Error(recorder, reader, http.StatusOK, tmpError)
	assert.Equal(t, http.StatusOK, recorder.Code)
	req := &models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(t, err)
	assert.Equal(t, req.Err, tmpError.Error())
}
