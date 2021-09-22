package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/store/teststore"
	"patreon/internal/models"
	"testing"
)

func TestMainHandler_HandleRegistration(t *testing.T) {
	s := teststore.New()
	TestOffLogger(t)
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
			recorder := httptest.NewRecorder()
			handler := NewMainHandler()
			handler.SetRouter(mux.NewRouter())
			handler.SetStore(s)
			handler.RegisterHandlers()

			b := bytes.Buffer{}

			err := json.NewEncoder(&b).Encode(test.data)
			assert.NoError(t, err)

			reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
			handler.ServeHTTP(recorder, reader)
			assert.Equal(t, test.expectedCode, recorder.Code)

		})
	}
}
func TestMainHandler_HandleLogin(t *testing.T) {
	u := models.TestUser(t)
	st := teststore.New()
	st.User().Create(u)
	var tests = []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"login":    u.Login,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid body",
			payload:      "hello, i'm incorrect body",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"login":    u.Login,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid login",
			payload: map[string]string{
				"login":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			handler := NewMainHandler()
			handler.SetRouter(mux.NewRouter())
			handler.SetStore(st)
			handler.RegisterHandlers()

			b := bytes.Buffer{}

			err := json.NewEncoder(&b).Encode(test.payload)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/login", &b)
			assert.NoError(t, err)

			handler.ServeHTTP(writer, req)
			assert.Equal(t, test.expectedCode, writer.Code)

		})
	}
}
