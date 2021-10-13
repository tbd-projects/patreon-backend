package creator_create_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/models"
	models_data "patreon/internal/app/models"
	"patreon/internal/app/repository"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreatorCreateTestSuite struct {
	handlers.SuiteHandler
	handler *CreatorCreateHandler
}

func (s *CreatorCreateTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewCreatorCreateHandler(s.Logger, s.MockSessionsManager, s.MockUserUsecase, s.MockCreatorUsecase)
}

func (s *CreatorCreateTestSuite) TestServeHTTP_Correct() {
	userID := int64(1)
	test := handlers.TestTable{
		Name:              "correct",
		Data:              userID,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.Data)

	require.NoError(s.T(), err)
	req, _ := http.NewRequest(http.MethodGet, "/creators", &b)
	vars := map[string]string{
		"id": strconv.Itoa(int(userID)),
	}
	creator := models_data.Creator{
		ID: userID, Avatar: "some", Nickname: "done"}
	reader := mux.SetURLVars(req, vars)
	s.MockCreatorUsecase.
		EXPECT().
		GetCreator(userID).
		Times(test.ExpectedMockTimes).
		Return(&creator, nil)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
	decoder := json.NewDecoder(recorder.Body)
	res := &models_data.Creator{}
	err = decoder.Decode(res)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), &creator, res)
}

func (s *CreatorCreateTestSuite) TestCreatorCreateHandler_POST_No_Params() {
	s.Tb = handlers.TestTable{
		Name:              "No url params",
		Data:              int64(-1),
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusBadRequest,
	}

	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(s.Tb.Data)

	require.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/creators", &b)

	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, &app.GeneralError{Err: repository.DefaultErrDB})
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *CreatorCreateTestSuite) TestCreatorCreateHandler_POST_Invalid_Body() {
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
	req, _ := http.NewRequest(http.MethodPost, "/creators", &b)
	vars := map[string]string{
		"id": strconv.Itoa(int(s.Tb.Data.(int64))),
	}
	reader := mux.SetURLVars(req, vars)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, &app.GeneralError{Err: repository.DefaultErrDB})
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}
func (s *CreatorCreateTestSuite) TestCreatorCreateHandler_POST_DB_Error() {
	s.Tb = handlers.TestTable{
		Name:              "Invalid request body",
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
	req, _ := http.NewRequest(http.MethodPost, "/creators", &b)
	vars := map[string]string{
		"id": strconv.Itoa(int(s.Tb.Data.(int64))),
	}
	reader := mux.SetURLVars(req, vars)
	s.MockUserUsecase.
		EXPECT().
		GetProfile(s.Tb.Data.(int64)).
		Times(s.Tb.ExpectedMockTimes).
		Return(nil, &app.GeneralError{Err: repository.DefaultErrDB})
	s.handler.POST(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
}

//func (s *CreatorCreateTestSuite) TestCreatorCreateHandler_POST_Create_Err() {
//	s.Tb = handlers.TestTable{
//		Name:              "CreateError in create usecase",
//		Data:              int64(1),
//		ExpectedMockTimes: 1,
//		ExpectedCode:      http.StatusBadRequest,
//	}
//	reqBody := models.RequestCreator{
//		Description: "description",
//		Category:    "category",
//	}
//	recorder := httptest.NewRecorder()
//
//	b := bytes.Buffer{}
//	err := json.NewEncoder(&b).Encode(reqBody)
//
//	require.NoError(s.T(), err)
//	req, _ := http.NewRequest(http.MethodPost, "/creators", &b)
//	vars := map[string]string{
//		"id": strconv.Itoa(int(s.Tb.Data.(int64))),
//	}
//	reader := mux.SetURLVars(req, vars)
//	s.MockUserUsecase.
//		EXPECT().
//		GetProfile(s.Tb.Data.(int64)).
//		Times(s.Tb.ExpectedMockTimes).
//		Return(nil, models_data.IncorrectCreatorCategory)
//	s.handler.POST(recorder, reader)
//	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
//}
func (s *CreatorCreateTestSuite) TestCreatorCreateHandler_POST_Correct() {
	userId := int64(1)
	s.Tb = handlers.TestTable{
		Name:              "Correct creator create",
		Data:              userId,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	reqBody := models.RequestCreator{
		Description: "description",
		Category:    "category",
	}
	creator := &models_data.Creator{
		ID:          userId,
		Category:    reqBody.Category,
		Description: reqBody.Description,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(reqBody)

	require.NoError(s.T(), err)
	req, _ := http.NewRequest(http.MethodPost, "/creators", &b)
	vars := map[string]string{
		"id": strconv.Itoa(int(s.Tb.Data.(int64))),
	}
	reader := mux.SetURLVars(req, vars)
	s.MockUserUsecase.
		EXPECT().
		Create(creator.ID).
		Times(s.Tb.ExpectedMockTimes).
		Return(creator.ID, nil)
	s.handler.POST(recorder, reader)
	decoder := json.NewDecoder(recorder.Body)
	var res int64
	err = decoder.Decode(res)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
	assert.Equal(s.T(), userId, res)

}
func TestCreatorCreateSuite(t *testing.T) {
	suite.Run(t, new(CreatorCreateTestSuite))
}
