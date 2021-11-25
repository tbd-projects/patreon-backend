package repository_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

type SuiteAttachesRepository struct {
	models.Suite
	repo *AttachesRepository
	data models.AttachWithoutLevel
}

func (s *SuiteAttachesRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewAttachesRepository(s.DB)
	s.data.Value = "asd"
	s.data.Type = "image"
	s.data.ID = 12
	s.data.PostId = 2
}

func (s *SuiteAttachesRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteAttachesRepository) TestAttachesRepository_getAndCheckDataTypeId() {
	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	id, err := s.repo.getAndCheckAttachTypeId(data.Type)
	assert.Equal(s.T(), id, int64(1))
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	id, err = s.repo.getAndCheckAttachTypeId(data.Type)
	assert.Equal(s.T(), id, int64(app.InvalidInt))
	assert.ErrorIs(s.T(), err, UnknownDataFormat)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.getAndCheckAttachTypeId(data.Type)
	assert.Equal(s.T(), id, int64(app.InvalidInt))
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteAttachesRepository) TestAttachesRepository_getDataType() {
	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(getDataTypeQuery)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(data.Type))
	id, err := s.repo.getAttachType(1)
	assert.Equal(s.T(), id, data.Type)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getDataTypeQuery)).
		WithArgs(1).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.getAttachType(1)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteAttachesRepository) TestAttachesRepository_Create() {
	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(1, data.Value, data.PostId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id"}).AddRow(data.ID))
	id, err := s.repo.Create(&data)
	assert.Equal(s.T(), id, data.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Create(&data)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(1, data.Value, data.PostId).
		WillReturnError(repository.DefaultErrDB)
	id, err = s.repo.Create(&data)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteAttachesRepository) TestAttachesRepository_Update() {
	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(updateQuery)).
		WithArgs(1, data.Value, data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"data_id"}).AddRow(data.ID))
	err := s.repo.Update(&data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(updateQuery)).
		WithArgs(1, data.Value, data.ID).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAndCheckDataTypeIdQuery)).
		WithArgs(data.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(updateQuery)).
		WithArgs(1, data.Value, data.ID).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Update(&data)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuiteAttachesRepository) TestAttachesRepository_Get() {
	data := s.data
	s.Mock.ExpectQuery(regexp.QuoteMeta(getQuery)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "data", "type"}).AddRow(data.PostId, data.Value, 1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(getDataTypeQuery)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(data.Type))
	res, err := s.repo.Get(data.ID)
	assert.Equal(s.T(), res, &data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getQuery)).
		WithArgs(data.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "data", "type"}).AddRow(data.PostId, data.Value, 1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(getDataTypeQuery)).
		WithArgs(1).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getQuery)).
		WithArgs(data.ID).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getQuery)).
		WithArgs(data.ID).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Get(data.ID)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuiteAttachesRepository) TestAttachesRepository_ExistsAttach() {
	attachId := int64(1)
	postId := int64(2)
	s.Mock.ExpectQuery(regexp.QuoteMeta(existsAttachQuery)).
		WithArgs(attachId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(postId))
	res, err := s.repo.ExistsAttach(attachId)
	assert.True(s.T(), res)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(existsAttachQuery)).
		WithArgs(attachId).
		WillReturnError(repository.DefaultErrDB)
	res, err = s.repo.ExistsAttach(attachId)
	assert.False(s.T(), res)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(existsAttachQuery)).
		WithArgs(attachId).
		WillReturnError(sql.ErrNoRows)
	res, err = s.repo.ExistsAttach(attachId)
	assert.False(s.T(), res)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuiteAttachesRepository) TestAttachesRepository_Delete() {
	attachId := int64(1)
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQuery)).
		WithArgs(attachId).
		WillReturnResult(driver.ResultNoRows)
	err := s.repo.Delete(attachId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQuery)).
		WithArgs(attachId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.Delete(attachId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuiteAttachesRepository) TestAttachesRepository_GetAttach() {
	data := s.data
	postId := data.PostId
	data.PostId = 0
	s.Mock.ExpectQuery(regexp.QuoteMeta(getAttachesQuery)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.Type, data.Value))
	res, err := s.repo.GetAttaches(postId)
	assert.Equal(s.T(), res[0], data)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAttachesQuery)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.Type, data.Value).
			RowError(0, repository.DefaultErrDB))
	_, err = s.repo.GetAttaches(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAttachesQuery)).
		WithArgs(postId).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.GetAttaches(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAttachesQuery)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"data_id", "type", "data"}).AddRow(data.ID, data.ID, data.Value))
	_, err = s.repo.GetAttaches(postId)
	assert.Error(s.T(), err)
}

func TestAttachesRepository(t *testing.T) {
	suite.Run(t, new(SuiteAttachesRepository))
}
