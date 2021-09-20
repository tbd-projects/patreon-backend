package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestOffLogger(t *testing.T) {
	t.Helper()
	logrus.SetOutput(ioutil.Discard)
}
