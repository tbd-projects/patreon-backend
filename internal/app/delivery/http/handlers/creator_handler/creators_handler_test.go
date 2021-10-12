package creator_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/delivery/http/handlers"
	models_data "patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CreatorTestSuite struct {
	handlers.SuiteHandler
	handler *CreatorHandler
}

func (s *CreatorTestSuite) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewCreatorHandler(s.Logger, s.MockSessionsManager, s.MockCreatorUsecase)
}

func (s *CreatorTestSuite) TestCreatorHandler_POST_Correct() {
	userID := int64(1)
	s.Tb = handlers.TestTable{
		Name:              "correct",
		Data:              &models.Creator{ID: userID, Avatar: "some", Nickname: "done"},
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
		Return([]models.Creator{*s.Tb.Data.(*models.Creator)}, nil)
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)
	req := &[]models_data.ResponseCreator{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &[]models_data.ResponseCreator{
		models_data.ToResponseCreator(*s.Tb.Data.(*models.Creator))})
}

func (s *CreatorTestSuite) TestCreatorHandler_POST_EmptyCreators() {
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
	s.handler.GET(recorder, reader)
	assert.Equal(s.T(), s.Tb.ExpectedCode, recorder.Code)

}
func TestCreatorHandler(t *testing.T) {
	suite.Run(t, new(CreatorTestSuite))
}
