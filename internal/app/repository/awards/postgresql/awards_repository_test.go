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

type SuiteAwardsRepository struct {
	models.Suite
	repo *AwardsRepository
}

func (s *SuiteAwardsRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewAwardsRepository(s.DB)
}

func (s *SuiteAwardsRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteAwardsRepository) TestAwardsRepository_checkUniqName() {
	query := "SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3"

	name := "sda"
	creatorId := int64(1)
	skipAwardsid := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	err := s.repo.checkUniqName(name, creatorId, skipAwardsid)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	err = s.repo.checkUniqName(name, creatorId, skipAwardsid)
	assert.Error(s.T(), err, NameAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnError(models.BDError)
	err = s.repo.checkUniqName(name, creatorId, skipAwardsid)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Create() {
	queryCheck := "SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3"
	query := `INSERT INTO awards (name, description, price, color, creator_id, cover) VALUES ($1, $2, $3, $4, $5. $6) 
				RETURNING awards_id`
	name := "sda"
	creatorId := int64(1)
	Id := int64(1)
	awards := models.Award{CreatorId: creatorId, Name: name}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, NotSkipAwards).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(Id))
	res, err := s.repo.Create(&awards)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res, Id)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, NotSkipAwards).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, NameAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, NotSkipAwards).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.CreatorId, app.DefaultImage).
		WillReturnError(models.BDError)
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Update() {
	queryCheck := "SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3"
	queryGetCreatorId := "SELECT creator_id from awards where awards.awards_id = $1"
	queryUpdate := "UPDATE awards SET name = $1, description = $2, price = $3, color = $4 WHERE awards_id = $5"
	name := "sda"
	creatorId := int64(1)
	Id := int64(1)
	awards := models.Award{CreatorId: creatorId, Name: name, ID: Id}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))
	err := s.repo.Update(&awards)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NotFound)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, NameAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnError(models.BDError)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(creatorId, name, Id).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(models.BDError))
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_CheckAwards() {
	query := `SELECT awards_id FROM awards where awards_id = $1`

	awardsId := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(awardsId))
	res, err := s.repo.CheckAwards(awardsId)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	res, err = s.repo.CheckAwards(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
	assert.False(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(awardsId).
		WillReturnError(sql.ErrNoRows)
	res, err = s.repo.CheckAwards(awardsId)
	assert.Error(s.T(), err, repository.NotFound)
	assert.False(s.T(), res)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_UpdateCover() {
	queryCheck := `SELECT creator_id from awards where awards.awards_id = $1`
	query := `UPDATE awards SET cover = $1 WHERE awards_id = $2`

	awardsId := int64(1)
	cover := "sad"
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(awardsId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}))
	err := s.repo.UpdateCover(awardsId, cover)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(awardsId).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NotFound)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(awardsId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, awardsId).
		WillReturnError(models.BDError)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(awardsId))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(models.BDError))
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_FindByName() {
	query := "SELECT count(*) as cnt from awards where creator_id = $1 and name = $2"

	creatorId := int64(1)
	name := "sad"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	res, err := s.repo.FindByName(creatorId, name)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	res, err = s.repo.FindByName(creatorId, name)
	assert.NoError(s.T(), err)
	assert.False(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId, name).
		WillReturnError(models.BDError)
	res, err = s.repo.FindByName(creatorId, name)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
	assert.False(s.T(), res)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_GetAwards() {
	query := `SELECT awards_id, name, description, price, color, cover from awards where awards.creator_id = $1`

	creatorId := int64(1)
	Id := int64(1)
	name := "sad"
	awards := models.Award{Name: name, ID: Id}

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price", "color", "cover"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.Cover))
	res, err := s.repo.GetAwards(creatorId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res[0], awards)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price", "color", "cover"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price, awards.Name, awards.Cover))
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price", "color", "cover"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.Cover).
			RowError(0, models.BDError))
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(creatorId).
		WillReturnError(models.BDError)
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_GetById() {
	query := "SELECT name, description, price, color, creator_id, cover FROM awards where awards_id = $1"

	creatorId := int64(1)
	Id := int64(2)
	name := "sad"
	awards := &models.Award{Name: name, ID: Id, CreatorId: creatorId}

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description", "price", "color", "creator_id", "cover"}).
			AddRow(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.CreatorId, awards.Cover))
	res, err := s.repo.GetByID(Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res, awards)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(Id).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetByID(Id)
	assert.Error(s.T(), err, repository.NotFound)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(Id).
		WillReturnError(models.BDError)
	_, err = s.repo.GetByID(Id)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Delete() {
	queryUpdate := "UPDATE posts SET type_awards = NULL where type_awards = $1"
	queryDelete := "DELETE FROM awards WHERE awards_id = $1"

	awardsId := int64(1)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	err := s.repo.Delete(awardsId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow().CloseError(models.BDError))
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow().CloseError(models.BDError))
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func TestAwardsRepository(t *testing.T) {
	suite.Run(t, new(SuiteAwardsRepository))
}
