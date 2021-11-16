package repository_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"github.com/zhashkevych/go-sqlxmock"
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

func (s *SuiteAwardsRepository) checkUniqCorrect(name string, creatorId int64, skipAwardsid int64, price int64) {
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryPrice)).
		WithArgs(creatorId, price, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
}

func (s *SuiteAwardsRepository) checkUniqError(name string, creatorId int64,
	skipAwardsid int64, _ int64, err error) {
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnError(err)
}

func (s *SuiteAwardsRepository) setLevelCorrect(awardsId int64, creatorId int64, price int64) {
	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertParent)).
		WithArgs(awardsId, price, creatorId).
		WillReturnResult(driver.RowsAffected(1))
	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertChild)).
		WithArgs(awardsId, price, creatorId).
		WillReturnResult(driver.RowsAffected(1))
}

func (s *SuiteAwardsRepository) setLevelError(awardsId int64, creatorId int64, price int64, err error) {
	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertParent)).
		WithArgs(awardsId, price, creatorId).
		WillReturnError(err)
}

func (s *SuiteAwardsRepository) deleteLevelCorrect(awardsId int64) {
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteLevelQuery)).
		WithArgs(awardsId).
		WillReturnResult(driver.RowsAffected(1))
}

func (s *SuiteAwardsRepository) deleteLevelError(awardsId int64, err error) {
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteLevelQuery)).
		WithArgs(awardsId).
		WillReturnError(err)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_checkUniqName() {
	name := "sda"
	creatorId := int64(1)
	skipAwardsid := int64(1)
	price := int64(3)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryPrice)).
		WithArgs(creatorId, price, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	err := s.repo.checkUniq(name, creatorId, skipAwardsid, price)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryPrice)).
		WithArgs(creatorId, price, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	err = s.repo.checkUniq(name, creatorId, skipAwardsid, price)
	assert.Error(s.T(), err, PriceAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	err = s.repo.checkUniq(name, creatorId, skipAwardsid, price)
	assert.Error(s.T(), err, NameAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnError(models.BDError)
	err = s.repo.checkUniq(name, creatorId, skipAwardsid, price)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryName)).
		WithArgs(creatorId, name, skipAwardsid).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	s.Mock.ExpectQuery(regexp.QuoteMeta(checkUniqQueryPrice)).
		WithArgs(creatorId, price, skipAwardsid).
		WillReturnError(repository.NewDBError(models.BDError))
	err = s.repo.checkUniq(name, creatorId, skipAwardsid, price)
	assert.Error(s.T(), err, PriceAlreadyExist)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_setLevel() {
	creatorId := int64(1)
	awardsId := int64(1)
	price := int64(3)

	s.Mock.ExpectBegin()
	trans, err := s.DB.Begin()
	require.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertParent)).
		WithArgs(awardsId, price, creatorId).
		WillReturnResult(driver.RowsAffected(1))
	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertChild)).
		WithArgs(awardsId, price, creatorId).
		WillReturnResult(driver.RowsAffected(1))
	err = s.repo.setLevel(trans, awardsId, price, creatorId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertParent)).
		WithArgs(awardsId, price, creatorId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.setLevel(trans, awardsId, price, creatorId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertParent)).
		WithArgs(awardsId, price, creatorId).
		WillReturnResult(driver.RowsAffected(1))
	s.Mock.ExpectExec(regexp.QuoteMeta(setLevelQueryInsertChild)).
		WithArgs(awardsId, price, creatorId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.setLevel(trans, awardsId, price, creatorId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectCommit()
	err = trans.Commit()
	require.NoError(s.T(), err)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_deleteLevel() {
	awardsId := int64(1)

	s.Mock.ExpectBegin()
	trans, err := s.DB.Begin()
	require.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(deleteLevelQuery)).
		WithArgs(awardsId).
		WillReturnResult(driver.RowsAffected(1))
	err = s.repo.deleteLevel(trans, awardsId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(deleteLevelQuery)).
		WithArgs(awardsId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.deleteLevel(trans, awardsId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectCommit()
	err = trans.Commit()
	require.NoError(s.T(), err)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Create() {
	name := "sda"
	creatorId := int64(1)
	Id := int64(1)
	awards := models.Award{CreatorId: creatorId, Name: name}

	s.checkUniqCorrect(name, creatorId, NotSkipAwards, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color),
			awards.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(Id))
	s.setLevelCorrect(Id, awards.CreatorId, awards.Price)
	s.Mock.ExpectCommit()
	res, err := s.repo.Create(&awards)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res, Id)

	s.checkUniqCorrect(name, creatorId, NotSkipAwards, awards.Price)
	s.Mock.ExpectBegin().WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.checkUniqCorrect(name, creatorId, NotSkipAwards, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color),
			awards.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(Id))
	s.setLevelCorrect(Id, awards.CreatorId, awards.Price)
	s.Mock.ExpectCommit().WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.checkUniqError(name, creatorId, NotSkipAwards, awards.Price, repository.DefaultErrDB)
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.checkUniqCorrect(name, creatorId, NotSkipAwards, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.CreatorId, app.DefaultImage).
		WillReturnError(models.BDError)
	s.Mock.ExpectRollback()
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.checkUniqCorrect(name, creatorId, NotSkipAwards, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(createQuery)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(Id))
	s.setLevelError(Id, awards.CreatorId, awards.Price, models.BDError)
	s.Mock.ExpectRollback()
	_, err = s.repo.Create(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Update() {
	name := "sda"
	creatorId := int64(1)
	Id := int64(1)
	awards := models.Award{CreatorId: creatorId, Name: name, ID: Id}

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(updateQueryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnResult(driver.RowsAffected(1))
	s.deleteLevelCorrect(Id)
	s.setLevelCorrect(Id, creatorId, awards.Price)
	s.Mock.ExpectCommit()
	err := s.repo.Update(&awards)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(updateQueryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnResult(driver.RowsAffected(1))
	s.deleteLevelCorrect(Id)
	s.setLevelCorrect(Id, creatorId, awards.Price)
	s.Mock.ExpectCommit().WillReturnError(repository.DefaultErrDB)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin().WillReturnError(repository.DefaultErrDB)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

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
	s.checkUniqError(name, creatorId, Id, awards.Price, NameAlreadyExist)
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, NameAlreadyExist)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(updateQueryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnError(models.BDError)
	s.Mock.ExpectRollback()
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(updateQueryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnResult(driver.RowsAffected(1))
	s.deleteLevelError(Id, models.BDError)
	s.Mock.ExpectRollback()
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	s.checkUniqCorrect(name, creatorId, Id, awards.Price)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(updateQueryUpdate)).
		WithArgs(awards.Name, awards.Description, awards.Price, convertRGBAToUint64(awards.Color), awards.ID).
		WillReturnResult(driver.RowsAffected(1))
	s.deleteLevelCorrect(Id)
	s.setLevelError(Id, creatorId, awards.Price, models.BDError)
	s.Mock.ExpectRollback()
	err = s.repo.Update(&awards)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_CheckAwards() {
	awardsId := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkAwardsQuery)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id"}).AddRow(awardsId))
	res, err := s.repo.CheckAwards(awardsId)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkAwardsQuery)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	res, err = s.repo.CheckAwards(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
	assert.False(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(checkAwardsQuery)).
		WithArgs(awardsId).
		WillReturnError(sql.ErrNoRows)
	res, err = s.repo.CheckAwards(awardsId)
	assert.Error(s.T(), err, repository.NotFound)
	assert.False(s.T(), res)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_UpdateCover() {
	awardsId := int64(1)
	cover := "sad"
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(awardsId))
	s.Mock.ExpectExec(regexp.QuoteMeta(updateCoverQuery)).
		WithArgs(cover, awardsId).WillReturnResult(driver.RowsAffected(1))
	err := s.repo.UpdateCover(awardsId, cover)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(awardsId).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NotFound)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryGetCreatorId)).
		WithArgs(awardsId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(awardsId))
	s.Mock.ExpectExec(regexp.QuoteMeta(updateCoverQuery)).
		WithArgs(cover, awardsId).
		WillReturnError(models.BDError)
	err = s.repo.UpdateCover(awardsId, cover)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_FindByName() {
	creatorId := int64(1)
	name := "sad"
	s.Mock.ExpectQuery(regexp.QuoteMeta(findByNameQuery)).
		WithArgs(creatorId, name).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	res, err := s.repo.FindByName(creatorId, name)
	assert.NoError(s.T(), err)
	assert.True(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(findByNameQuery)).
		WithArgs(creatorId, name).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	res, err = s.repo.FindByName(creatorId, name)
	assert.NoError(s.T(), err)
	assert.False(s.T(), res)

	s.Mock.ExpectQuery(regexp.QuoteMeta(findByNameQuery)).
		WithArgs(creatorId, name).
		WillReturnError(models.BDError)
	res, err = s.repo.FindByName(creatorId, name)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
	assert.False(s.T(), res)
}

func (s *SuiteAwardsRepository) TestAwardsRepository_GetAwards() {
	creatorId := int64(1)
	Id := int64(1)
	name := "sad"
	awards := models.Award{Name: name, ID: Id}

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAwardsQuery)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price",
			"color", "cover", "child_id"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price,
				convertRGBAToUint64(awards.Color), awards.Cover, awards.ChildAward))
	res, err := s.repo.GetAwards(creatorId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res[0], awards)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAwardsQuery)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price",
			"color", "cover", "child_id"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price,
					awards.Cover, awards.Cover, awards.ChildAward))
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAwardsQuery)).
		WithArgs(creatorId).
		WillReturnRows(sqlmock.NewRows([]string{"awards_id", "name", "description", "price",
			"color", "cover", "child_id"}).
			AddRow(awards.ID, awards.Name, awards.Description, awards.Price,
				convertRGBAToUint64(awards.Color), awards.Cover, awards.ChildAward).
			RowError(0, models.BDError))
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(getAwardsQuery)).
		WithArgs(creatorId).
		WillReturnError(models.BDError)
	_, err = s.repo.GetAwards(creatorId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_GetById() {
	creatorId := int64(1)
	Id := int64(2)
	name := "sad"
	awards := &models.Award{Name: name, ID: Id, CreatorId: creatorId}

	s.Mock.ExpectQuery(regexp.QuoteMeta(getByIdQuery)).
		WithArgs(Id).
		WillReturnRows(sqlmock.NewRows([]string{"name", "description", "price", "color", "creator_id", "cover", "child_id"}).
			AddRow(awards.Name, awards.Description, awards.Price,
				convertRGBAToUint64(awards.Color), awards.CreatorId, awards.Cover, awards.ChildAward))
	res, err := s.repo.GetByID(Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res, awards)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getByIdQuery)).
		WithArgs(Id).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetByID(Id)
	assert.Error(s.T(), err, repository.NotFound)

	s.Mock.ExpectQuery(regexp.QuoteMeta(getByIdQuery)).
		WithArgs(Id).
		WillReturnError(models.BDError)
	_, err = s.repo.GetByID(Id)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuiteAwardsRepository) TestAwardsRepository_Delete() {
	awardsId := int64(1)
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQueryUpdate)).
		WithArgs(awardsId).
		WillReturnResult(driver.RowsAffected(1))
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQueryDelete)).
		WithArgs(awardsId).
		WillReturnResult(driver.RowsAffected(1))
	err := s.repo.Delete(awardsId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQueryUpdate)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQueryUpdate)).
		WithArgs(awardsId).
		WillReturnResult(driver.RowsAffected(1))
	s.Mock.ExpectExec(regexp.QuoteMeta(deleteQueryDelete)).
		WithArgs(awardsId).
		WillReturnError(models.BDError)
	err = s.repo.Delete(awardsId)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func TestAwardsRepository(t *testing.T) {
	suite.Run(t, new(SuiteAwardsRepository))
}
