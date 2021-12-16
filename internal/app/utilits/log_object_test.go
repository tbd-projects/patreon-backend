package utilits

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestLogObject(t *testing.T) {
	defer func(t *testing.T) {
		err := recover()
		require.Equal(t, err, nil)
	}(t)

	log := &logrus.Logger{}
	object := NewLogObject(log)

	b := bytes.Buffer{}

	reader, err := http.NewRequest(http.MethodPost, "/register", &b)
	require.NoError(t, err)

	object.BaseLog()
	assert.Equal(t, object.BaseLog(), log)

	entry := object.Log(reader)
	assert.Equal(t, entry.Data["urls"].(*url.URL), reader.URL)

	entry = object.Log(nil)
	assert.Equal(t, entry.Data["type"].(string), "base_log")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", log.WithField("tamp", "da"))
	reader, err = http.NewRequestWithContext(ctx, http.MethodPost, "/register", &b)
	entry = object.Log(reader)
	assert.Equal(t, entry.Data["tamp"].(string), "da")
}
