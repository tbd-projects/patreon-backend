package creator_create_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	models_data "patreon/internal/app/models"
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
	s.handler = NewCreatorCreateHandler(s.Logger, s.Router, s.Cors, s.MockSessionsManager, s.MockUserUsecase, s.MockCreatorUsecase)
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

func TestCreatorCreateSuite(t *testing.T) {
	suite.Run(t, new(CreatorCreateTestSuite))
}
