package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_awards "patreon/internal/app/repository/awards"

	"github.com/pkg/errors"
)

const NotSkipAwards = -1

type AwardsRepository struct {
	store *sql.DB
}

var _ = repository_awards.Repository(&AwardsRepository{})

func NewAwardsRepository(st *sql.DB) *AwardsRepository {
	return &AwardsRepository{
		store: st,
	}
}

// Create Errors:
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) checkUniqName(name string, creatorId int64, skipAwardsid int64) error {
	query := "SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3"
	count := 0
	if err := repo.store.QueryRow(query, creatorId, name, skipAwardsid).Scan(&count); err != nil {
		return repository.NewDBError(err)
	}

	if count != 0 {
		return NameAlreadyExist
	}

	return nil
}

// Create Errors:
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Create(aw *models.Award) (int64, error) {
	if err := repo.checkUniqName(aw.Name, aw.CreatorId, NotSkipAwards); err != nil {
		return -1, err
	}

	query := `INSERT INTO awards (name, description, price, color, creator_id, cover) VALUES ($1, $2, $3, $4, $5. $6, $7) 
				RETURNING awards_id`

	if err := repo.store.QueryRow(query, aw.Name, aw.Description, aw.Price, convertRGBAToUint64(aw.Color),
		aw.CreatorId, app.DefaultImage).
		Scan(&aw.ID); err != nil {
		return -1, repository.NewDBError(err)
	}
	return aw.ID, nil
}

// GetByID Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetByID(awardsID int64) (*models.Award, error) {
	query := "SELECT name, description, price, color, creator_id, cover FROM awards where awards_id = $1"
	aw := &models.Award{ID: awardsID}
	var clr uint64
	if err := repo.store.QueryRow(query,
		awardsID).Scan(&aw.Name, &aw.Description, &aw.Price, &clr, &aw.CreatorId, &aw.Cover); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}
	aw.Color = convertUint64ToRGBA(clr)
	return aw, nil
}

// CheckAwards Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) CheckAwards(awardsID int64) (bool, error) {
	query := `SELECT awards_id FROM awards where awards_id = $1`

	if err := repo.store.QueryRow(query, awardsID).Scan(&awardsID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.NotFound
		}
		return false, repository.NewDBError(err)
	}

	return true, nil
}

// Update Errors:
//		repository.NotFound
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Update(aw *models.Award) error {
	queryGetCreatorId := "SELECT creator_id from awards where awards.awards_id = $1"
	queryUpdate := "UPDATE awards SET name = $1, description = $2, price = $3, color = $4 WHERE awards_id = $5"
	creatorId := int64(0)

	if err := repo.store.QueryRow(queryGetCreatorId, aw.ID).Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if err := repo.checkUniqName(aw.Name, creatorId, aw.ID); err != nil {
		return err
	}

	row, err := repo.store.Query(queryUpdate, aw.Name, aw.Description, aw.Price, convertRGBAToUint64(aw.Color), aw.ID)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// GetAwards Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetAwards(creatorId int64) ([]models.Award, error) {
	query := `SELECT awards_id, name, description, price, color, cover from awards where awards.creator_id = $1`
	var res []models.Award

	rows, err := repo.store.Query(query, creatorId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		var awards models.Award
		var clr uint64
		if err = rows.Scan(&awards.ID, &awards.Name, &awards.Description, &awards.Price, &clr, &awards.Cover); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		awards.Color = convertUint64ToRGBA(clr)
		res = append(res, awards)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Delete(awardsId int64) error {
	queryUpdate := "UPDATE posts SET type_awards = NULL where type_awards = $1"
	queryDelete := "DELETE FROM awards WHERE awards_id = $1"

	row, err := repo.store.Query(queryUpdate, awardsId)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}

	row, err = repo.store.Query(queryDelete, awardsId)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// FindByName Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) FindByName(creatorID int64, awardName string) (bool, error) {
	query := "SELECT count(*) as cnt from awards where creator_id = $1 and name = $2"
	cnt := 0
	res := repo.store.QueryRow(query, creatorID, awardName).Scan(&cnt)
	if res != nil {
		return false, repository.NewDBError(res)
	}
	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// UpdateCover Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) UpdateCover(awardsId int64, cover string) error {
	queryCheck := `SELECT creator_id from awards where awards.awards_id = $1`
	query := `UPDATE awards SET cover = $1 WHERE awards_id = $2`
	creatorId := int64(0)

	if err := repo.store.QueryRow(queryCheck, awardsId).Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	row, err := repo.store.Query(query, cover, awardsId)
	if err != nil {
		return repository.NewDBError(err)
	}

	if err = row.Close(); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}
