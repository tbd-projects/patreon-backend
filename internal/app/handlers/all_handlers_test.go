package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app"
	"patreon/internal/app/sessions/mocks"
	"patreon/internal/app/store"
	mock_store "patreon/internal/app/store/mocks"
	"testing"
)

type TestTable struct {
	name              string
	data              interface{}
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

type SuiteTestBaseHandler struct {
	suite.Suite
	mock                  *gomock.Controller
	mockUserRepository    *mock_store.MockUserRepository
	mockCreatorRepository *mock_store.MockCreatorRepository
	mockSessionsManager   *mocks.MockSessionsManager
	store                 store.Store
	test                  TestTable
}

func (s *SuiteTestBaseHandler) SetupSuite() {
	s.mock = gomock.NewController(s.T())
	s.mockUserRepository = mock_store.NewMockUserRepository(s.mock)
	s.mockCreatorRepository = mock_store.NewMockCreatorRepository(s.mock)
	s.mockSessionsManager = mocks.NewMockSessionsManager(s.mock)

	s.store = NewStore(s.mockUserRepository, s.mockCreatorRepository)

	s.test = TestTable{}
}

func (s *SuiteTestBaseHandler) TearDownSuite() {
	s.mock.Finish()
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
	suite.Run(t, new(LoginTestSuite))
	suite.Run(t, new(ProfileTestSuite))
	suite.Run(t, new(RegisterTestSuite))
	suite.Run(t, new(CreatorTestSuite))
	t.Run("Join run", func(t *testing.T) {
		defer func() {
			err := recover()
			require.Equal(t, err, nil)
		}()

		router := mux.NewRouter()
		handler := NewMainHandler()
		handler.SetRouter(router)
		registerHandler := NewRegisterHandler()
		loginHandler := NewLoginHandler()
		profileHandler := NewProfileHandler()
		logoutHandler := NewLogoutHandler()
		creatorHandler := NewCreatorHandler()
		creatorCreateHandler := NewCreatorCreateHandler()

		creatorHandler.JoinHandlers([]app.Joinable{
			creatorCreateHandler,
		})

		handler.JoinHandlers([]app.Joinable{
			registerHandler,
			loginHandler,
			profileHandler,
			logoutHandler,
			creatorHandler,
		})
	})
}
