package likes_handler

import (
	"patreon/internal/app/delivery/http/handlers"
)

type LikesIdTestSuite struct {
	handlers.SuiteHandler
	handler *LikesHandler
}

func (s *LikesIdTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	//	s.handler = NewLikesHandler(s.Logger, s.Router, s.Cors, s.MockAwardsUsecase, s.MockSessionsManager)
}

/*
func (s *LikesIdTestSuite) TestCreatorIdHandler_POST_No_Params() {
	s.Tb = handlers.TestTable{
		Name:              "No url params",
		Data:              int64(-1),
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusInternalServerError,
	}

	reqBody := models.RequestCreator{
		Description: "description",
		Category:    "category",
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(reqBody)

	require.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/creators", &b)

	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes)
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *LikesIdTestSuite) TestCreatorIdHandler_POST_Invalid_Body() {
	s.Tb = handlers.TestTable{
		Name:              "Invalid request body",
		Data:              int64(1),
		ExpectedMockTimes: 0,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.Tb.Data)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", s.Tb.Data.(int64))
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/creators", &b)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, &app.GeneralError{Err: repository.DefaultErrDB})
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *LikesIdTestSuite) TestCreatorIdHandler_POST_DB_Error() {
	s.Tb = handlers.TestTable{
		Name:              "Invalid request body",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	reqBody := models.RequestCreator{
		Description: "description",
		Category:    "category",
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(reqBody)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", s.Tb.Data.(int64))
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/creators", &b)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, &app.GeneralError{Err: repository.DefaultErrDB})
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *LikesIdTestSuite) TestCreatorIdHandler_POST_Create_Err() {
	s.Tb = handlers.TestTable{
		Name:              "CreateError in create usecase",
		Data:              int64(1),
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	reqBody := models.RequestCreator{
		Description: "description",
		Category:    "category",
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(reqBody)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", s.Tb.Data.(int64))
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/creators", &b)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, models_data.IncorrectCreatorCategory)
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func (s *LikesIdTestSuite) TestCreatorIdHandler_POST_Correct() {
	userId := int64(1)
	user := models_data.TestUser()
	user.ID = userId
	user.Nickname = "nickname"
	s.Tb = handlers.TestTable{
		Name:              "Correct creator create",
		Data:              userId,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusCreated,
	}
	reqBody := models.RequestCreator{
		Description: "description",
		Category:    "category",
	}
	creator := &models_data.Creator{
		ID:          user.ID,
		Nickname:    user.Nickname,
		Category:    reqBody.Category,
		Description: reqBody.Description,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(reqBody)

	require.NoError(s.T(), err)
	ctx := context.WithValue(context.Background(), "user_id", s.Tb.Data.(int64))
	reader, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/creators", &b)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(creator.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(user, nil)
	s.MockCreatorUsecase.
		EXPECT().
		Create(newCreatorWithFieldMatcher(creator)).
		Times(s.Tb.ExpectedMockTimes).
		Return(creator.ID, nil)
	s.handler.ServeHTTP(recorder, reader)
	decoder := json.NewDecoder(recorder.Body)
	var res interface{}
	err = decoder.Decode(&res)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
	assert.Equal(s.T(), int(userId), int(res.(float64)))

}

type creatorWithFieldMatcher struct{ creator *models_data.Creator }

func newCreatorWithFieldMatcher(creator *models_data.Creator) gomock.Matcher {
	return &creatorWithFieldMatcher{creator}
}

func (match *creatorWithFieldMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case *models_data.Creator:
		return x.(*models_data.Creator).ID == match.creator.ID && x.(*models_data.Creator).Nickname == match.creator.Nickname &&
			x.(*models_data.Creator).Category == match.creator.Category && x.(*models_data.Creator).Description == match.creator.Description
	default:
		return false
	}
}

func (match *creatorWithFieldMatcher) String() string {
	return fmt.Sprintf("Creator: %s", match.creator.String())
}

func (s *LikesIdTestSuite) TestAwardsHandler_GET_Correct() {
	userID := int64(1)
	s.Tb = handlers.TestTable{
		Name:              "correct",
		Data:              &models_data.Creator{ID: userID, Avatar: "some", Nickname: "done"},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.Tb.Data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodGet, "/creators", &b)

	s.MockCreatorUsecase.
		EXPECT().
		GetCreators().
		Times(s.Tb.ExpectedMockTimes).
		Return([]models_data.Creator{*s.Tb.Data.(*models_data.Creator)}, nil)
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
	req := &[]models.ResponseCreator{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &[]models.ResponseCreator{
		models.ToResponseCreator(*s.Tb.Data.(*models_data.Creator))})
}

func (s *LikesIdTestSuite) TestAwardsHandler_GET_EmptyCreators() {
	s.Tb = handlers.TestTable{
		Name:              "creators is empty",
		Data:              nil,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.Tb.Data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodGet, "/creators", &b)

	s.MockCreatorUsecase.
		EXPECT().
		GetCreators().
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, nil)
	s.handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

func TestAwardsHandler(t *testing.T) {
	suite.Run(t, new(LikesIdTestSuite))
}
*/
