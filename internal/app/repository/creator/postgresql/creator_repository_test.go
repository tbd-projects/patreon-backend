package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	rp "patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/zhashkevych/go-sqlxmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type SuiteCreatorRepository struct {
	models.Suite
	repo *CreatorRepository
}

func (s *SuiteCreatorRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewCreatorRepository(s.DB)
}

func (s *SuiteCreatorRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteCreatorRepository) TestCreatorRepository_Create() {
	queryCategory := `SELECT category_id FROM creator_category WHERE name = $1`

	query := `INSERT INTO creator_profile (creator_id, category,
		description, avatar, cover) VALUES ($1, $2, $3, $4, $5)
		RETURNING creator_id
	`
	cr := models.TestCreator()

	cr.ID = 1
	categoryId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).
		WithArgs(cr.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(categoryId))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID, categoryId, cr.Description, app.DefaultImage, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(cr.ID))))
	id, err := s.repo.Create(cr)
	assert.Equal(s.T(), id, cr.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).
		WithArgs(cr.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(categoryId))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID, categoryId, cr.Description, app.DefaultImage, app.DefaultImage).WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).
		WithArgs(cr.Category).
		WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).
		WithArgs(cr.Category).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), IncorrectCategory, err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreator() {
	cr := models.TestCreatorWithAwards()
	cr.ID = 1
	userId := int64(1)
	expected := *cr

	var awardsId sql.NullInt64
	awardsId.Valid = true
	awardsId.Int64 = cr.AwardsId
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname", "awards_id"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname, awardsId))

	get, err := s.repo.GetCreator(userId, expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)


	awardsId.Valid = false
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname", "awards_id"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname, awardsId))
	expected.AwardsId = rp.NoAwards
	get, err = s.repo.GetCreator(userId, expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).WillReturnError(sql.ErrNoRows)

	_, err = s.repo.GetCreator(userId, expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NotFound, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreator)).
		WithArgs(userId, cr.ID).WillReturnError(models.BDError)

	_, err = s.repo.GetCreator(userId, expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreators_AllUsersCreators() {
	creators := models.TestCreators()

	preapareRows := sqlmock.NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"})

	for index, cr := range creators {
		cr.ID = int64(index)
		creators[index] = cr
		preapareRows.AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname)
	}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(preapareRows)

	get, err := s.repo.GetCreators()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), creators, get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(preapareRows.RowError(0, models.BDError))

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(""))

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountGetCreators)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreatorGetCreators)).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateAvatar() {
	avatar := "d"
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	err := s.repo.UpdateAvatar(creatorId, avatar)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).WillReturnError(sql.ErrNoRows)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), app.UnknownError, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAvatar)).WithArgs(avatar, creatorId).WillReturnError(models.BDError)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateCover() {
	cover := "d"
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateCover)).WithArgs(cover, creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	err := s.repo.UpdateCover(creatorId, cover)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateCover)).WithArgs(cover, creatorId).WillReturnError(sql.ErrNoRows)

	err = s.repo.UpdateCover(creatorId, cover)
	assert.Error(s.T(), app.UnknownError, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdateCover)).WithArgs(cover, creatorId).WillReturnError(models.BDError)

	err = s.repo.UpdateCover(creatorId, cover)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_ExistCreator() {
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	check, err := s.repo.ExistsCreator(creatorId)
	assert.NoError(s.T(), err)
	assert.True(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).WillReturnError(sql.ErrNoRows)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), app.UnknownError, err)
	assert.False(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryExistsCreator)).WithArgs(creatorId).WillReturnError(models.BDError)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
	assert.False(s.T(), check)
}

func TestCreatorRepository(t *testing.T) {
	suite.Run(t, new(SuiteCreatorRepository))
}
