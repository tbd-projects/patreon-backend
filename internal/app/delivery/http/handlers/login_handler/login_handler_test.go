package login_handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/mailru/easyjson"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/models"
	model_data "patreon/internal/app/models"
	session_models "patreon/internal/microservices/auth/sessions/models"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type LoginTestSuite struct {
	handlers.SuiteHandler
	handler *LoginHandler
}

func (s *LoginTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewLoginHandler(s.Logger, s.MockSessionsManager, s.MockUserUsecase)
}

func (s *LoginTestSuite) TestLoginHandler_POST_EmptyBody() {
	s.Tb = handlers.TestTable{
		Name:              "Empty body in request",
		Data:              &http_models.RequestLogin{},
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *LoginTestSuite) TestLoginHandler_POST_InvalidBody() {
	s.Tb = handlers.TestTable{
		Name:              "Invalid body",
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	data := http_models.ResponsePost{
		Title: "nickname",
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *LoginTestSuite) TestLoginHandler_POST_UserNotFound() {
	s.Tb = handlers.TestTable{
		Name: "User not found in db",
		Data: http_models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnauthorized,
	}

	recorder := httptest.NewRecorder()
	expectedId := int64(-1)
	s.MockUserUsecase.EXPECT().
		Check(s.Tb.Data.(http_models.RequestLogin).Login,
			s.Tb.Data.(http_models.RequestLogin).Password).
		Times(s.Tb.ExpectedMockTimes).
		Return(expectedId, model_data.IncorrectEmailOrPassword)

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *LoginTestSuite) TestLoginHandler_POST_SessionError() {
	s.Tb = handlers.TestTable{
		Name: "Create Session Error",
		Data: http_models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	recorder := httptest.NewRecorder()

	user := model_data.User{
		ID:       1,
		Login:    s.Tb.Data.(http_models.RequestLogin).Login,
		Password: s.Tb.Data.(http_models.RequestLogin).Password,
	}
	err := user.Encrypt()
	assert.NoError(s.T(), err)
	s.MockUserUsecase.EXPECT().
		Check(s.Tb.Data.(http_models.RequestLogin).Login,
			s.Tb.Data.(http_models.RequestLogin).Password).
		Times(s.Tb.ExpectedMockTimes).
		Return(user.ID, nil)

	s.MockSessionsManager.EXPECT().
		Create(context.Background(), int64(user.ID)).
		Times(s.Tb.ExpectedMockTimes).
		Return(session_models.Result{
			UserID: -1,
		},
			errors.New("error"))

	b := bytes.Buffer{}
	_, err = easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *LoginTestSuite) TestLoginHandler_POST_Ok() {
	s.Tb = handlers.TestTable{
		Name: "Invalid body",
		Data: http_models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	recorder := httptest.NewRecorder()

	user := model_data.User{
		ID:       1,
		Login:    s.Tb.Data.(http_models.RequestLogin).Login,
		Password: s.Tb.Data.(http_models.RequestLogin).Password,
	}
	err := user.Encrypt()
	assert.NoError(s.T(), err)
	s.MockUserUsecase.EXPECT().
		Check(s.Tb.Data.(http_models.RequestLogin).Login,
			s.Tb.Data.(http_models.RequestLogin).Password).
		Times(s.Tb.ExpectedMockTimes).
		Return(user.ID, nil)

	s.MockSessionsManager.EXPECT().
		Create(context.Background(), int64(user.ID)).
		Times(s.Tb.ExpectedMockTimes).
		Return(session_models.Result{UserID: 1, UniqID: "123"}, nil)

	b := bytes.Buffer{}
	_, err = easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	s.handler.POST(recorder, reader)

	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func TestLoginHandler(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
