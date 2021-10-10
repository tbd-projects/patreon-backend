package handlers

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMainHandler_HandleRegistration(t *testing.T) {
	//s := teststore.New()
	logrus.SetOutput(ioutil.Discard)
	var tests = []struct {
		name         string
		data         interface{}
		expectedCode int
	}{
		{
			name: "valid",
			data: map[string]string{
				"login":    "yandex@mail.ru",
				"password": "qwerty",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "incorrect login",
			data: map[string]string{
				"login":    "cat",
				"password": "qwerty",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty password",
			data: map[string]string{
				"login":    "cat1998",
				"password": "",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request format",
			data: map[string]string{
				"login":    "cat",
				"email":    "login@email.com",
				"password": "",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty request body",
			data:         map[string]string{},
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//recorder := httptest.NewRecorder()
			//handler := NewMainHandler()
			//handler.SetRouter(mux.NewRouter())
			//handler.SetStore(s)
			//handler.RegisterHandlers()

			//b := bytes.Buffer{}

			//err := json.NewEncoder(&b).Encode(test.data)
			//assert.NoError(t, err)
			//assert.NoError(t, nil)

			//reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
			//handler.ServeHTTP(recorder, reader)
			//assert.Equal(t, test.expectedCode, recorder.Code)

		})
	}
}
