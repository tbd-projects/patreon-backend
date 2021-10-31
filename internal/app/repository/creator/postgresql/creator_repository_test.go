package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
		WithArgs(cr.ID, categoryId, cr.Description, cr.Avatar, cr.Cover).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(cr.ID))))
	id, err := s.repo.Create(cr)
	assert.Equal(s.T(), id, cr.ID)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCategory)).
		WithArgs(cr.Category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(strconv.Itoa(int(categoryId))))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID, categoryId, cr.Description, cr.Avatar, cr.Cover).WillReturnError(models.BDError)
	_, err = s.repo.Create(cr)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreator() {
	query := `SELECT creator_id, cc.name, description, creator_profile.avatar, cover, usr.nickname 
			FROM creator_profile JOIN users AS usr ON usr.users_id = creator_profile.creator_id 
			JOIN creator_category As cc ON creator_profile.category = cc.category_id 
			where creator_id=$1`
	cr := models.TestCreator()
	cr.ID = 1
	expected := *cr

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"}).
			AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname))

	get, err := s.repo.GetCreator(expected.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expected, *get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID).WillReturnError(sql.ErrNoRows)

	_, err = s.repo.GetCreator(expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NotFound, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cr.ID).WillReturnError(models.BDError)

	_, err = s.repo.GetCreator(expected.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_GetCreators_AllUsersCreators() {
	queryCount := `SELECT count(*) from creator_profile`
	queryCreator := `SELECT creator_id, cc.name, description, creator_profile.avatar, cover, usr.nickname 
					FROM creator_profile JOIN users AS usr ON usr.users_id = creator_profile.creator_id
					JOIN creator_category As cc ON creator_profile.category = cc.category_id`

	creators := models.TestCreators()

	preapareRows := sqlmock.NewRows([]string{"id", "category", "description", "avatar", "cover", "nickname"})

	for index, cr := range creators {
		cr.ID = int64(index)
		creators[index] = cr
		preapareRows.AddRow(strconv.Itoa(int(cr.ID)), cr.Category, cr.Description, cr.Avatar, cr.Cover, cr.Nickname)
	}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreator)).
		WillReturnRows(preapareRows)

	get, err := s.repo.GetCreators()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), creators, get)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(strconv.Itoa(len(creators))))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCreator)).WillReturnError(models.BDError)

	_, err = s.repo.GetCreators()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateAvatar() {
	query := `UPDATE creator_profile SET avatar = $1 WHERE creator_id = $2 RETURNING creator_id`

	avatar := "d"
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(avatar, creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	err := s.repo.UpdateAvatar(creatorId, avatar)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(avatar, creatorId).WillReturnError(sql.ErrNoRows)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), app.UnknownError, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(avatar, creatorId).WillReturnError(models.BDError)

	err = s.repo.UpdateAvatar(creatorId, avatar)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_UpdateCover() {
	query := `UPDATE creator_profile SET cover = $1 WHERE creator_id = $2 RETURNING creator_id`

	cover := "d"
	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(cover, creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	err := s.repo.UpdateCover(creatorId, cover)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(cover, creatorId).WillReturnError(sql.ErrNoRows)

	err = s.repo.UpdateCover(creatorId, cover)
	assert.Error(s.T(), app.UnknownError, err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(cover, creatorId).WillReturnError(models.BDError)

	err = s.repo.UpdateCover(creatorId, cover)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
}

func (s *SuiteCreatorRepository) TestCreatorRepository_ExistCreator() {
	query := `SELECT creator_id from creator_profile where creator_id=$1`

	creatorId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))

	check, err := s.repo.ExistsCreator(creatorId)
	assert.NoError(s.T(), err)
	assert.True(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(creatorId).WillReturnError(sql.ErrNoRows)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), app.UnknownError, err)
	assert.False(s.T(), check)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(creatorId).WillReturnError(models.BDError)

	check, err = s.repo.ExistsCreator(creatorId)
	assert.Error(s.T(), repository.NewDBError(models.BDError), err)
	assert.False(s.T(), check)
}

func TestCreatorRepository(t *testing.T) {
	suite.Run(t, new(SuiteCreatorRepository))
}
