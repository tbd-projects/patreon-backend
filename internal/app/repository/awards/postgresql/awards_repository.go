package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_awards "patreon/internal/app/repository/awards"
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
	count := 0
	err := repo.store.QueryRow(
		"SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3",
		creatorId, name, skipAwardsid).Scan(&count)
	if err != nil {
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
func (repo *AwardsRepository) Create(aw *models.Awards) (int64, error) {
	if err := repo.checkUniqName(aw.Name, aw.CreatorId, NotSkipAwards); err != nil {
		return -1, err
	}

	if err := repo.store.QueryRow("INSERT INTO awards (name, "+
		"description, price, color, creator_id) VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING awards_id", aw.Name, aw.Description, aw.Price, convertRGBAToUint64(aw.Color), aw.CreatorId).
		Scan(&aw.ID); err != nil {
		return -1, repository.NewDBError(err)
	}
	return aw.ID, nil
}

// GetByID Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetByID(awardsID int64) (*models.Awards, error) {
	aw := &models.Awards{ID: awardsID}
	var clr uint64
	if err := repo.store.QueryRow("SELECT name, description, price, color, creator_id FROM awards where awards_id = $1",
		aw.ID).Scan(&aw.Name, &aw.Description, &aw.Price, &clr, &aw.CreatorId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}
	aw.Color = convertUint64ToRGBA(clr)
	return aw, nil
}

// Update Errors:
//		repository.NotFound
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Update(aw *models.Awards) error {
	creatorId := int64(0)

	if err := repo.store.QueryRow(
		"SELECT creator_id from awards where awards.awards_id = $1", aw.ID).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if err := repo.checkUniqName(aw.Name, creatorId, aw.ID); err != nil {
		return err
	}

	if _, err := repo.store.Query("UPDATE awards SET name = $1, description = $2, price = $3, color = $4 WHERE awards_id = $5",
		aw.Name, aw.Description, aw.Price, convertRGBAToUint64(aw.Color), aw.ID); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// GetAwards Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetAwards(creatorId int64) ([]models.Awards, error) {
	var res []models.Awards

	rows, err := repo.store.Query(
		"SELECT awards_id, name, description, price, color from awards where awards.creator_id = $1", creatorId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var awards models.Awards
		var clr uint64
		if err = rows.Scan(&awards.ID, &awards.Name, &awards.Description, &awards.Price, &clr); err != nil {
			return nil, repository.NewDBError(err)
		}
		awards.Color = convertUint64ToRGBA(clr)
		res = append(res, awards)
		i++

		if err = rows.Err(); err != nil {
			return nil, repository.NewDBError(err)
		}
	}
	if err = rows.Close(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Delete(awardsId int64) error {
	if _, err := repo.store.Query("UPDATE posts SET type_awards = NULL where type_awards = $1",
		awardsId); err != nil {
		return repository.NewDBError(err)
	}

	if _, err := repo.store.Query("DELETE FROM awards WHERE awards_id = $1", awardsId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
