package repository_postgresql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

type SuitePostsDataRepository struct {
	models.Suite
	repo *PostsDataRepository
	data models.PostData
}

func (s *SuitePostsDataRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewPostsDataRepository(s.DB)
	s.data.Data = "asd"
	s.data.Type = "image"
	s.data.ID = 12
	s.data.PostId = 2
}

func (s *SuitePostsDataRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}


func (s *SuitePostsDataRepository) TestPostsDataRepository_getAndCheckDataTypeId() {
	query := `SELECT posts_type_id FROM posts_type WHERE type = $1`

	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	id, err := s.repo.getAndCheckDataTypeId(data.Type)
	assert.Equal(s.T(), id, int64(1))
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	id, err = s.repo.getAndCheckDataTypeId(data.Type)
	assert.Equal(s.T(), id, int64(app.InvalidInt))
	assert.ErrorIs(s.T(), err, UnknownDataFormat)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.Type).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.getAndCheckDataTypeId(data.Type)
	assert.Equal(s.T(), id, int64(app.InvalidInt))
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuitePostsDataRepository) TestPostsDataRepository_getDataType() {
	query := `SELECT type FROM posts_type WHERE posts_type_id = $1`

	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(data.Type))
	id, err := s.repo.getDataType(1)
	assert.Equal(s.T(), id, data.Type)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.getDataType(1)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}



func (s *SuitePostsDataRepository) TestPostsDataRepository_Create() {
	queryCheck := `SELECT posts_type_id FROM posts_type WHERE type = $1`
	query := `INSERT INTO posts_data (type, data, post_id) VALUES ($1, $2, $3) 
		RETURNING data_id`

	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, data.Data, data.PostId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id"}).AddRow(data.ID))
	id, err := s.repo.Create(&data)
	assert.Equal(s.T(), id, data.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Create(&data)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, data.Data, data.PostId).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.Create(&data)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}


func (s *SuitePostsDataRepository) TestPostsDataRepository_Update() {
	queryCheck := `SELECT posts_type_id FROM posts_type WHERE type = $1`
	query := `UPDATE posts_data SET type = $1, data = $2 WHERE data_id = $3 RETURNING data_id`

	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, data.Data, data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"data_id"}).AddRow(data.ID))
	err := s.repo.Update(&data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, data.Data, data.ID).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, data.Data, data.ID).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsDataRepository) TestPostsDataRepository_Get() {
	queryCheck := `SELECT type FROM posts_type WHERE posts_type_id = $1`
	query := `SELECT post_id, data, type FROM posts_data WHERE data_id = $1`

	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "data", "type"}).AddRow(data.PostId, data.Data, 1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(data.Type))
	res, err := s.repo.Get(data.ID)
	assert.Equal(s.T(), res, &data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "data", "type"}).AddRow(data.PostId, data.Data, 1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(1).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.ID).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(data.ID).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsDataRepository) TestPostsDataRepository_ExistsData() {
	query := `SELECT post_id FROM posts_data WHERE data_id = $1`

	dataId := int64(1)
	postId := int64(2)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(dataId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(postId))
	res, err := s.repo.ExistsData(dataId)
	assert.True(s.T(), res)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnError(repository.DefaultErrDB)
	res, err = s.repo.ExistsData(dataId)
	assert.False(s.T(), res)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	res, err = s.repo.ExistsData(dataId)
	assert.False(s.T(), res)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsDataRepository) TestPostsDataRepository_Delete() {
	query := `DELETE FROM posts_data WHERE data_id = $1`

	dataId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(dataId).
		WillReturnRows(sqlmock.NewRows([]string{}))
	err := s.repo.Delete(dataId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))
	err = s.repo.Delete(dataId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Delete(dataId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuitePostsDataRepository) TestPostsDataRepository_GetData() {
	query := `SELECT data_id, pst.type, data FROM posts_data JOIN posts_type AS pst 
    			ON (pst.posts_type_id = posts_data.type) WHERE post_id = $1`

	data := s.data
	postId := data.PostId
	data.PostId = 0
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.Type, data.Data))
	res, err := s.repo.GetData(postId)
	assert.Equal(s.T(), res[0], data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.Type, data.Data).
			RowError(0, repository.DefaultErrDB))
	_, err = s.repo.GetData(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.GetData(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.ID, data.Data))
	_, err = s.repo.GetData(postId)
	assert.Error(s.T(), err)
}

func TestPostsDataRepository(t *testing.T) {
	suite.Run(t, new(SuitePostsDataRepository))
}

