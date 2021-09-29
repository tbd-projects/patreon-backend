package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/store"
	mock_store "patreon/internal/app/store/mocks"
	"patreon/internal/models"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

type TestTable struct {
	name              string
	data              *models.RequestLogin
	expectedMockTimes int
	expectedCode      int
}
type Store struct {
	userRepository    store.UserRepository
	creatorRepository store.CreatorRepository
}

func NewStore(userRep store.UserRepository,
	creatorRep store.CreatorRepository) *Store {
	return &Store{userRep, creatorRep}
}

func (st *Store) User() store.UserRepository {
	return st.userRepository
}
func (st *Store) Creator() store.CreatorRepository {
	return st.creatorRepository
}

type SuiteTestStore struct {
	suite.Suite
	mock                  *gomock.Controller
	mockUserRepository    *mock_store.MockUserRepository
	mockCreatorRepository *mock_store.MockCreatorRepository
	store                 store.Store
	test                  TestTable
}

func (s *SuiteTestStore) SetupSuite() {
	s.mock = gomock.NewController(s.T())
	s.mockUserRepository = mock_store.NewMockUserRepository(s.mock)
	s.mockCreatorRepository = mock_store.NewMockCreatorRepository(s.mock)

	s.store = NewStore(s.mockUserRepository, s.mockCreatorRepository)

	s.test = TestTable{}
}

func (s *SuiteTestStore) TearDownSuite() {
	s.mock.Finish()
}

func TestTestStore(t *testing.T) {
	suite.Run(t, new(SuiteTestStore))
}

func (s *SuiteTestStore) TestLoginHandler_ServeHTTP_EmptyBody() {
	s.test = TestTable{
		name:              "Invalid body",
		data:              &models.RequestLogin{},
		expectedMockTimes: 0,
		expectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *SuiteTestStore) TestLoginHandler_ServeHTTP_InvalidBody() {
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
	handler := NewLoginHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
func (s *SuiteTestStore) TestLoginHandler_ServeHTTP_UserNotFound() {
	s.test = TestTable{
		name: "Invalid body",
		data: &models.RequestLogin{
			Login:    "dmitriy",
			Password: "mail.ru",
		},
		expectedMockTimes: 1,
		expectedCode:      http.StatusUnauthorized,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)

	handler.SetStore(s.store)
	s.mockUserRepository.EXPECT().
		FindByLogin(s.test.data.Login).
		Times(s.test.expectedMockTimes).
		Return(nil, store.NotFound)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.test.expectedCode, recorder.Code)
}
