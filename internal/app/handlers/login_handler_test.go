package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app"
	session_models "patreon/internal/app/sessions/models"
	"patreon/internal/app/sessions/sessions_manager"

	"patreon/internal/app/store"
	"patreon/internal/models"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

type LoginTestSuite struct {
	SuiteTestBaseHandler
}

func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_EmptyBody() {
	s.test = TestTable{
		name:              "Empty body in request",
		data:              &models.RequestLogin{},
		expectedMockTimes: 0,
		expectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()
	dataStorage := &app.DataStorage{
		Store:          s.store,
		SessionManager: s.mockSessionsManager,
	}
	handler := NewLoginHandler(dataStorage)
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_InvalidBody() {
	s.test = TestTable{
		name:              "Invalid body",
		expectedMockTimes: 0,
		expectedCode:      http.StatusUnprocessableEntity,
	}
	data := struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}{
		Nickname: "nickname",
		Password: "password",
	}
	recorder := httptest.NewRecorder()
	dataStorage := &app.DataStorage{
		Store:          s.store,
		SessionManager: s.mockSessionsManager,
	}
	handler := NewLoginHandler(dataStorage)
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_UserNotFound() {
	s.test = TestTable{
		name: "User not found in db",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusUnauthorized,
	}

	recorder := httptest.NewRecorder()
	dataStorage := &app.DataStorage{
		Store:          s.store,
		SessionManager: s.mockSessionsManager,
	}
	handler := NewLoginHandler(dataStorage)
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.(models.RequestLogin).Login).
		Times(s.test.expectedMockTimes).
		Return(nil, store.NotFound)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_UserNoAuthorized() {
	s.test = TestTable{
		name: "Not authorized user",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusInternalServerError,
	}
	recorder := httptest.NewRecorder()
	dataStorage := &app.DataStorage{
		Store:          s.store,
		SessionManager: s.mockSessionsManager,
	}
	handler := NewLoginHandler(dataStorage)
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	user := models.User{
		ID:       1,
		Login:    s.test.data.(models.RequestLogin).Login,
		Password: s.test.data.(models.RequestLogin).Password,
	}
	err := user.BeforeCreate()
	assert.NoError(s.T(), err)
	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.(models.RequestLogin).Login).
		Times(s.test.expectedMockTimes).
		Return(&models.User{ID: user.ID, Login: user.Login, EncryptedPassword: user.EncryptedPassword}, nil)

	s.mockSessionsManager.EXPECT().
		Create(int64(user.ID)).
		Times(s.test.expectedMockTimes).
		Return(session_models.Result{UserID: sessions_manager.UnknownUser},
			errors.New("error"))

	b := bytes.Buffer{}
	err = json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_Ok() {
	s.test = TestTable{
		name: "Invalid body",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusOK,
	}
	recorder := httptest.NewRecorder()
	dataStorage := &app.DataStorage{
		Store:          s.store,
		SessionManager: s.mockSessionsManager,
	}
	handler := NewLoginHandler(dataStorage)
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	user := models.User{
		ID:       1,
		Login:    s.test.data.(models.RequestLogin).Login,
		Password: s.test.data.(models.RequestLogin).Password,
	}
	err := user.BeforeCreate()
	assert.NoError(s.T(), err)
	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.(models.RequestLogin).Login).
		Times(s.test.expectedMockTimes).
		Return(&models.User{ID: user.ID, Login: user.Login, EncryptedPassword: user.EncryptedPassword}, nil)

	s.mockSessionsManager.EXPECT().
		Create(int64(user.ID)).
		Times(s.test.expectedMockTimes).
		Return(session_models.Result{UserID: 1, UniqID: "123"}, nil)

	b := bytes.Buffer{}
	err = json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
