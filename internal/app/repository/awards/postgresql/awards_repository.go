package repository_postgresql

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type AwardsRepository struct {
	store *sql.DB
}

func NewAwardsRepository(st *sql.DB) *AwardsRepository {
	return &AwardsRepository{
		store: st,
	}
}

// Create Errors:
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) checkUniqName(name string, creatorId int64) error {
	count := 0
	err := repo.store.QueryRow(
		"SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2", creatorId, name).Scan(&count)
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
	if err := repo.checkUniqName(aw.Name, aw.CreatorId); err != nil {
		return -1, err
	}

	if err := repo.store.QueryRow("INSERT INTO awards (name, "+
		"description, price, creator_id) VALUES ($1, $2, $3, $4)"+
		"RETURNING awards_id", aw.Name, aw.Description, aw.Price, aw.CreatorId).Scan(&aw.ID); err != nil {
		return -1, repository.NewDBError(err)
	}
	return aw.ID, nil
}

// UpdateName Errors:
//		repository.NotFound
//		repository_postgresql.NameAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) UpdateName(awardsId int64, name string) error {
	creatorId := int64(0)

	if err := repo.store.QueryRow(
		"SELECT creator_id from awards where awards.awards_id = $1", awardsId).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if err := repo.checkUniqName(name, creatorId); err != nil {
		return err
	}

	if _, err := repo.store.Query("UPDATE awards SET name = $1 WHERE awards_id = $2",
		name, awardsId); err != nil {
		return repository.NewDBError(err)
	}
	return nil
}

// UpdatePriceDescription Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) UpdatePriceDescription(awardsId int64, price int64, description string) error {
	if _, err := repo.store.Query("UPDATE awards SET description = $1, price = $2 WHERE awards_id = $3",
		description, price, awardsId); err != nil {
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
		"SELECT awards_id, name, description, price from awards where awards.creator_id = $1", creatorId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var awards models.Awards
		if err = rows.Scan(&awards.ID, &awards.Name, &awards.Description, &awards.Price); err != nil {
			return nil, repository.NewDBError(err)
		}
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
func (repo *AwardsRepository) Delete(awardsId int64) (*models.Creator, error) {
	creator := &models.Creator{}

	if _, err := repo.store.Query("UPDATE posts SET type_awards = NULL where type_awards = $1",
		awardsId); err != nil {
		return nil, repository.NewDBError(err)
	}

	if _, err := repo.store.Query("DELETE FROM awards WHERE awards_id = $1", awardsId); err != nil {
		return nil, repository.NewDBError(err)
	}

	return creator, nil
}
