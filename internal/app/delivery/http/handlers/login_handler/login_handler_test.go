package login_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	models2 "patreon/internal/app/repository/models"
	session_models "patreon/internal/app/sessions/models"
	"patreon/internal/app/sessions/sessions_manager"

	"patreon/internal/app/store"
	"patreon/internal/models"

	"github.com/stretchr/testify/assert"
)

type LoginTestSuite struct {
	handlers.SuiteTestBaseHandler
}

func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_EmptyBody() {
	s.test = handlers.TestTable{
		name:              "Empty body in request",
		data:              &models.RequestLogin{},
		expectedMockTimes: 0,
		expectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler(s.logger, s.dataStorage)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_InvalidBody() {
	s.test = handlers.TestTable{
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
	handler := NewLoginHandler(s.logger, s.dataStorage)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_ServeHTTP_UserNotFound() {
	s.test = handlers.TestTable{
		name: "User not found in db",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusUnauthorized,
	}

	recorder := httptest.NewRecorder()
	handler := NewLoginHandler(s.logger, s.dataStorage)

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
	s.test = handlers.TestTable{
		name: "Not authorized user",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusInternalServerError,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler(s.logger, s.dataStorage)

	user := models2.User{
		ID:       1,
		Login:    s.test.data.(models.RequestLogin).Login,
		Password: s.test.data.(models.RequestLogin).Password,
	}
	err := user.Encrypt()
	assert.NoError(s.T(), err)
	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.(models.RequestLogin).Login).
		Times(s.test.expectedMockTimes).
		Return(&models2.User{ID: user.ID, Login: user.Login, EncryptedPassword: user.EncryptedPassword}, nil)

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
	s.test = handlers.TestTable{
		name: "Invalid body",
		data: models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusOK,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler(s.logger, s.dataStorage)

	user := models2.User{
		ID:       1,
		Login:    s.test.data.(models.RequestLogin).Login,
		Password: s.test.data.(models.RequestLogin).Password,
	}
	err := user.Encrypt()
	assert.NoError(s.T(), err)
	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.(models.RequestLogin).Login).
		Times(s.test.expectedMockTimes).
		Return(&models2.User{ID: user.ID, Login: user.Login, EncryptedPassword: user.EncryptedPassword}, nil)

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
