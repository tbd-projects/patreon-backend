package search_creators_handler

import (
	"patreon/internal/app/delivery/http/handlers"
)

type SearchCreatorsTestSuite struct {
	handlers.SuiteHandler
	handler *SearchCreatorsHandler
}

func (s *SearchCreatorsTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewCreatorHandler(s.Logger, s.MockSessionsManager, s.MockCreatorUsecase)
}
/*
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

func (s *SearchCreatorsTestSuite) TestCreatorHandler_GET_Correct() {
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
	req := &[]http_models.ResponseCreator{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &[]http_models.ResponseCreator{
		http_models.ToResponseCreator(*s.Tb.Data.(*models_data.Creator))})
}

func (s *SearchCreatorsTestSuite) TestCreatorHandler_GET_EmptyCreators() {
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

func TestCreatorHandler(t *testing.T) {
	suite.Run(t, new(SearchCreatorsTestSuite))
}
*/