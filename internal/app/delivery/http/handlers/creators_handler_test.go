package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"patreon/internal/app/store"
	"patreon/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CreatorTestSuite struct {
	SuiteTestBaseHandler
}

func (s *CreatorTestSuite) TestServeHTTP_Correct() {
	userID := int64(1)
	test := TestTable{
		name:              "correct",
		data:              &models.Creator{ID: int(userID), Avatar: "some", Nickname: "done"},
		expectedMockTimes: 1,
		expectedCode:      http.StatusOK,
	}

	recorder := httptest.NewRecorder()
	handler := NewCreatorHandler(s.logger, s.dataStorage)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodGet, "/creators", &b)

	s.mockCreatorRepository.
		EXPECT().
		GetCreators().
		Times(test.expectedMockTimes).
		Return([]models.Creator{*test.data.(*models.Creator)}, nil)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)

	req := &[]models.ResponseCreator{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(req)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), req, &[]models.ResponseCreator{models.ToResponseCreator(*test.data.(*models.Creator))})
}

func (s *CreatorTestSuite) TestServeHTTP_WitDBError() {
	test := TestTable{
		name:              "with db error",
		data:              nil,
		expectedMockTimes: 1,
		expectedCode:      http.StatusServiceUnavailable,
	}

	recorder := httptest.NewRecorder()
	handler := NewCreatorHandler(s.logger, s.dataStorage)

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(test.data)

	require.NoError(s.T(), err)
	reader, _ := http.NewRequest(http.MethodGet, "/creators", &b)

	s.mockCreatorRepository.
		EXPECT().
		GetCreators().
		Times(test.expectedMockTimes).
		Return(nil, store.NotFound)
	handler.ServeHTTP(recorder, reader)
	assert.Equal(s.T(), test.expectedCode, recorder.Code)

}
