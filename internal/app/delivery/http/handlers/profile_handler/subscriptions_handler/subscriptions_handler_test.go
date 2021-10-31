package subscriptions_handler

import (
	"patreon/internal/app/delivery/http/handlers"
)

type SubscriptionsTestSuite struct {
	handlers.SuiteHandler
	handler *SubscriptionsHandler
}

//func (s *SubscriptionsTestSuite) SetupSuite() {
//	s.SuiteHandler.SetupSuite()
//	s.handler = NewSubscriptionsHandler(s.Logger, s.Router, s.Cors,
//		s.MockSessionsManager, s.MockSubscribersUsecase)
//}
//func (s *SubscriptionsTestSuite) TestSubscriptionsHandler_GET_OK_NoOneCreators() {
//	userID := int64(1)
//	expectedCreators := responseModels.SubscriptionsUserResponse{}
//	s.Tb = handlers.TestTable{
//		Name:              "correct test",
//		Data:              "",
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusOK,
//	}
//	writer := httptest.NewRecorder()
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(s.Tb.Data)
//	require.NoError(s.T(), err)
//	ctx := context.WithValue(context.Background(), "user_id", userID)
//	reader, err := http.NewRequestWithContext(ctx, http.MethodGet, "/user/subscriptions", &b)
//
//	s.MockSubscribersUsecase.EXPECT().
//		GetCreators(userID).
//		Times(1).
//		Return(expectedCreators.Creators, nil)
//
//	s.handler.GET(writer, reader)
//
//	assert.Equal(s.T(), s.Tb.ExpectedCode, writer.Code)
//
//	resCreators := responseModels.SubscriptionsUserResponse{}
//	decoder := json.NewDecoder(writer.Body)
//	err = decoder.Decode(&resCreators)
//
//	require.NoError(s.T(), err)
//	assert.Equal(s.T(), expectedCreators, resCreators)
//
//}
//func (s *SubscriptionsTestSuite) TestSubscriptionsHandler_GET_OK_OneCreator() {
//	userID := int64(1)
//	expectedCreators := responseModels.SubscriptionsUserResponse{
//		Creators: []int64{1},
//	}
//	s.Tb = handlers.TestTable{
//		Name:              "correct test",
//		Data:              "",
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusOK,
//	}
//	writer := httptest.NewRecorder()
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(s.Tb.Data)
//	require.NoError(s.T(), err)
//	ctx := context.WithValue(context.Background(), "user_id", userID)
//	reader, err := http.NewRequestWithContext(ctx, http.MethodGet, "/user/subscriptions", &b)
//
//	s.MockSubscribersUsecase.EXPECT().
//		GetCreators(userID).
//		Times(1).
//		Return(expectedCreators.Creators, nil)
//
//	s.handler.GET(writer, reader)
//
//	assert.Equal(s.T(), s.Tb.ExpectedCode, writer.Code)
//
//	resCreators := responseModels.SubscriptionsUserResponse{}
//	decoder := json.NewDecoder(writer.Body)
//	err = decoder.Decode(&resCreators)
//
//	require.NoError(s.T(), err)
//	assert.Equal(s.T(), expectedCreators, resCreators)
//
//}
//func (s *SubscriptionsTestSuite) TestSubscriptionsHandler_GET_OK_FewCreators() {
//	userID := int64(1)
//	expectedCreators := responseModels.SubscriptionsUserResponse{
//		Creators: []int64{1, 2, 3},
//	}
//	s.Tb = handlers.TestTable{
//		Name:              "correct test",
//		Data:              "",
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusOK,
//	}
//	writer := httptest.NewRecorder()
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(s.Tb.Data)
//	require.NoError(s.T(), err)
//	ctx := context.WithValue(context.Background(), "user_id", userID)
//	reader, err := http.NewRequestWithContext(ctx, http.MethodGet, "/user/subscriptions", &b)
//
//	s.MockSubscribersUsecase.EXPECT().
//		GetCreators(userID).
//		Times(1).
//		Return(expectedCreators.Creators, nil)
//
//	s.handler.GET(writer, reader)
//
//	assert.Equal(s.T(), s.Tb.ExpectedCode, writer.Code)
//
//	resCreators := responseModels.SubscriptionsUserResponse{}
//	decoder := json.NewDecoder(writer.Body)
//	err = decoder.Decode(&resCreators)
//
//	require.NoError(s.T(), err)
//	assert.Equal(s.T(), expectedCreators, resCreators)
//
//}
//func (s *SubscriptionsTestSuite) TestSubscriptionsHandler_GET_ContextError() {
//	s.Tb = handlers.TestTable{
//		Name:              "can not get user_id",
//		Data:              "",
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusInternalServerError,
//	}
//	writer := httptest.NewRecorder()
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(s.Tb.Data)
//	require.NoError(s.T(), err)
//	reader, err := http.NewRequest(http.MethodGet, "/user/subscriptions", &b)
//
//	s.handler.GET(writer, reader)
//
//	assert.Equal(s.T(), s.Tb.ExpectedCode, writer.Code)
//
//}
//func (s *SubscriptionsTestSuite) TestSubscriptionsHandler_GET_UsecaseError() {
//	userID := int64(1)
//	s.Tb = handlers.TestTable{
//		Name:              "usecase return error",
//		Data:              "",
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusInternalServerError,
//	}
//	writer := httptest.NewRecorder()
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(s.Tb.Data)
//	require.NoError(s.T(), err)
//	ctx := context.WithValue(context.Background(), "user_id", userID)
//	reader, err := http.NewRequestWithContext(ctx, http.MethodGet, "/user/subscriptions", &b)
//
//	s.MockSubscribersUsecase.EXPECT().
//		GetCreators(userID).
//		Times(1).
//		Return(nil, repository.NewDBError(db_models.BDError))
//
//	s.handler.GET(writer, reader)
//
//	assert.Equal(s.T(), s.Tb.ExpectedCode, writer.Code)
//
//	expErr := responseModels.ErrResponse{
//		Err: handler_errors.BDError.Error(),
//	}
//	resErr := responseModels.ErrResponse{}
//	decoder := json.NewDecoder(writer.Body)
//	err = decoder.Decode(&resErr)
//
//	require.Equal(s.T(), expErr, resErr)
//}
//func TestSubscriptionsHandler(t *testing.T) {
//	suite.Run(t, new(SubscriptionsTestSuite))
//}
