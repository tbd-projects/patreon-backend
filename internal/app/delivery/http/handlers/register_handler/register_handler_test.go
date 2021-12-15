package register_handler

import (
	"bytes"
	"fmt"
	"github.com/mailru/easyjson"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/models"
	models_data "patreon/internal/app/models"
	repository_user "patreon/internal/app/repository/user/postgresql"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type RegisterTestSuite struct {
	handlers.SuiteHandler
	handler *RegisterHandler
}

func (s *RegisterTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewRegisterHandler(s.Logger, s.MockSessionsManager, s.MockUserUsecase)
}

func (s *RegisterTestSuite) TestRegisterHandler_POST_EmptyBody() {
	s.Tb = handlers.TestTable{
		Name:              "Empty body from request",
		Data:              &http_models.RequestRegistration{},
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}

	reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *RegisterTestSuite) TestRegisterHandler_POST_InvalidBody() {
	s.Tb = handlers.TestTable{
		Name:              "Invalid body",
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	Data := http_models.ResponsePost {
		Title:    "nickname",
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *RegisterTestSuite) TestRegisterHandler_POST_UserAlreadyExist() {
	s.Tb = handlers.TestTable{
		Name: "User exist in database",
		Data: http_models.RequestRegistration{
			Login:    "dmitriy",
			Nickname: "linux1998",
			Password: "mail.ru",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusConflict,
	}

	recorder := httptest.NewRecorder()

	req := s.Tb.Data.(http_models.RequestRegistration)
	user := &models_data.User{
		Login:    req.Login,
		Nickname: req.Nickname,
		Password: req.Password,
	}
	expId := int64(-1)
	s.MockUserUsecase.EXPECT().
		Create(user).
		Times(s.Tb.ExpectedMockTimes).
		Return(expId, repository_user.LoginAlreadyExist)

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *RegisterTestSuite) TestRegisterHandler_POST_SmallPassword() {
	s.Tb = handlers.TestTable{
		Name: "Small password in request",
		Data: http_models.RequestRegistration{
			Login:    "dmitriy",
			Nickname: "linux1998",
			Password: "mail",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	req := s.Tb.Data.(http_models.RequestRegistration)
	user := &models_data.User{
		Login:    req.Login,
		Nickname: req.Nickname,
		Password: req.Password,
	}
	expId := int64(-1)
	s.MockUserUsecase.EXPECT().
		Create(user).
		Times(s.Tb.ExpectedMockTimes).
		Return(expId, models_data.IncorrectEmailOrPassword)

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(s.Tb.Data, &b)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

type userWithPasswordMatcher struct{ user *models_data.User }

func newUserWithPasswordMatcher(user *models_data.User) gomock.Matcher {
	return &userWithPasswordMatcher{user}
}

func (match *userWithPasswordMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case *models_data.User:
		return x.(*models_data.User).ID == match.user.ID && x.(*models_data.User).Login == match.user.Login &&
			x.(*models_data.User).Avatar == match.user.Avatar && x.(*models_data.User).Nickname == match.user.Nickname &&
			x.(*models_data.User).Password == match.user.Password && match.user.ComparePassword(x.(*models_data.User).Password)
	default:
		return false
	}
}

func (match *userWithPasswordMatcher) String() string {
	return fmt.Sprintf("User: %s", match.user.String())
}

func (s *RegisterTestSuite) TestRegisterHandler_POST_CreateSuccess() {
	s.Tb = handlers.TestTable{
		Name: "Success create user",
		Data: http_models.RequestRegistration{
			Login:    "dmitriy",
			Password: "mail.ru",
			Nickname: "linux1998",
		},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusCreated,
	}
	recorder := httptest.NewRecorder()

	req := s.Tb.Data.(http_models.RequestRegistration)
	user := &models_data.User{
		Login:    req.Login,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	b := bytes.Buffer{}
	_, err := easyjson.MarshalToWriter(s.Tb.Data, &b)
	assert.NoError(s.T(), err)

	assert.NoError(s.T(), user.Encrypt())

	s.MockUserUsecase.EXPECT().
		Create(newUserWithPasswordMatcher(user)).
		Times(1).
		Return(user.ID, nil)

	reader, _ := http.NewRequest(http.MethodPost, "/register", &b)
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func TestRegisterHandler(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
