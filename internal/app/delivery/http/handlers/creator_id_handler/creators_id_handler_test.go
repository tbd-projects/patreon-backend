package creator_id_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	"patreon/internal/app/delivery/http/models"
	models_data "patreon/internal/app/models"
	usecase_creator "patreon/internal/app/usecase/creator"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreatorCreateTestSuite struct {
	handlers.SuiteHandler
	handler *CreatorIdHandler
}

func (s *CreatorCreateTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewCreatorIdHandler(s.Logger, s.MockSessionsManager, s.MockUserUsecase, s.MockCreatorUsecase)
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
		"creator_id": strconv.Itoa(int(userID)),
	}
	creator := models_data.CreatorWithAwards{
		ID: userID, Avatar: "some", Nickname: "done"}
	reader := mux.SetURLVars(req, vars)
	s.MockCreatorUsecase.
		EXPECT().
		GetCreator(userID, usecase_creator.NoUser).
		Times(test.ExpectedMockTimes).
		Return(&creator, nil)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), test.ExpectedCode, recorder.Code)
	decoder := json.NewDecoder(recorder.Body)
	res := &models.ResponseCreatorWithAwards{}
	err = decoder.Decode(res)
	assert.NoError(s.T(), err)
	expected := models.ToResponseCreatorWithAwards(creator)
	assert.Equal(s.T(), &expected, res)
}

func TestCreatorCreateSuite(t *testing.T) {
	suite.Run(t, new(CreatorCreateTestSuite))
}
