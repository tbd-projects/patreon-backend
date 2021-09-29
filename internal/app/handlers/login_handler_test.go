package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	store "patreon/internal/app/store"
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
	data              *models.User
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
	tests                 []TestTable
}

func (s *SuiteTestStore) SetupSuite() {
	s.mock = gomock.NewController(s.T())
	s.mockUserRepository = mock_store.NewMockUserRepository(s.mock)
	s.mockCreatorRepository = mock_store.NewMockCreatorRepository(s.mock)

	s.store = NewStore(s.mockUserRepository, s.mockCreatorRepository)

	s.tests = []TestTable{}
}

func (s *SuiteTestStore) TearDownSuite() {
	s.mock.Finish()
}

func TestTestStore(t *testing.T) {
	suite.Run(t, new(SuiteTestStore))
}

func (s *SuiteTestStore) TestLoginHandler_ServeHTTP_EmptyBody() {
	type TestTable struct {
		name              string
		data              *models.User
		expectedMockTimes int
		expectedCode      int
	}
	test := TestTable{
		name:              "Invalid body",
		data:              &models.User{},
		expectedMockTimes: 0,
		expectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()
	handler := NewLoginHandler()
	logger := logrus.New()
	str := bytes.Buffer{}
	logger.SetOutput(&str)

	handler.SetLogger(logger)

	//store := new(Store)
	//handler.SetStore(store)

	//s.mockUserRepository.EXPECT().
	//	Create(s.data).
	//	Times(test.expectedMockTimes).
	//	Do()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/login", &b)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)
}
