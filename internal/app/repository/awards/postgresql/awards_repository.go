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

const (
	checkUniqQueryName  = "SELECT count(*) from awards where awards.creator_id = $1 and awards.name = $2 and awards.awards_id != $3"
	checkUniqQueryPrice = "SELECT count(*) from awards where awards.creator_id = $1 and awards.price = $2 and awards.awards_id != $3"

	setLevelQueryInsertParent = `INSERT INTO parents_awards (awards_id, parent_id) SELECT $1, awards_id 
					FROM awards WHERE price > $2 AND creator_id = $3 ORDER BY price`
	setLevelQueryInsertChild = `INSERT INTO parents_awards (awards_id, parent_id) SELECT awards_id, $1 
					FROM awards WHERE price < $2 AND creator_id = $3 ORDER BY price`

	deleteLevelQuery = `DELETE FROM parents_awards WHERE awards_id = $1 OR parent_id = $1`

	createQuery = `INSERT INTO awards (name, description, price, color, creator_id, cover) VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING awards_id`

	queryGetCreatorId = "SELECT creator_id from awards where awards.awards_id = $1"
	updateQueryUpdate = "UPDATE awards SET name = $1, description = $2, price = $3, color = $4 WHERE awards_id = $5"

	updateCoverQuery = `UPDATE awards SET cover = $1 WHERE awards_id = $2`

	getByIdQuery = `with frist_child AS (
						SELECT parents_awards.awards_id as award_id, parents_awards.parent_id as parent_id, price FROM awards a
						JOIN parents_awards ON parents_awards.awards_id = a.awards_id and parents_awards.parent_id = $1
						ORDER BY price DESC LIMIT 1
					)
					SELECT aw.name, aw.description, aw.price, aw.color, aw.creator_id, aw.cover, ch.award_id as child_id FROM awards AS aw
    				LEFT JOIN frist_child as ch on ch.parent_id = $1 WHERE aw.awards_id = $1`

	checkAwardsQuery = `SELECT awards_id FROM awards where awards_id = $1`

	getAwardsQuery = `with frist_child AS (
						SELECT parents_awards.parent_id as parent_id, max(price) as mx_price
						FROM awards a
								 JOIN parents_awards ON parents_awards.awards_id = a.awards_id
						WHERE a.creator_id = $1 GROUP BY parent_id
					), child_with_price AS (
						SELECT parents_awards.awards_id as award_id, parents_awards.parent_id as parent_id, price
						FROM awards a
								 JOIN parents_awards ON parents_awards.awards_id = a.awards_id
						WHERE a.creator_id = $1
					)
					SELECT aw.awards_id, aw.name, aw.description, aw.price, aw.color, aw.cover, pa.award_id as child_id
					FROM awards AS aw
							 LEFT JOIN frist_child as ch on ch.parent_id = aw.awards_id
							 LEFT JOIN child_with_price pa on ch.parent_id = pa.parent_id and ch.mx_price = pa.price
					WHERE aw.creator_id = $1;`

	deleteQueryUpdate = "UPDATE posts SET type_awards = NULL where type_awards = $1"
	deleteQueryDelete = "DELETE FROM awards WHERE awards_id = $1"

	findByNameQuery = "SELECT count(*) as cnt from awards where creator_id = $1 and name = $2"
)

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
func (repo *AwardsRepository) checkUniq(name string, creatorId int64, skipAwardsid int64, price int64) error {
	count := 0
	if err := repo.store.QueryRow(checkUniqQueryName, creatorId, name, skipAwardsid).Scan(&count); err != nil {
		return repository.NewDBError(err)
	}

	if count != 0 {
		return NameAlreadyExist
	}

	count = 0
	if err := repo.store.QueryRow(checkUniqQueryPrice, creatorId, price, skipAwardsid).Scan(&count); err != nil {
		return repository.NewDBError(err)
	}

	if count != 0 {
		return PriceAlreadyExist
	}

	return nil
}

// Create Errors:
//		repository_postgresql.NameAlreadyExist
//		repository_postgresql.PriceAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Create(aw *models.Award) (int64, error) {
	if err := repo.checkUniq(aw.Name, aw.CreatorId, NotSkipAwards, aw.Price); err != nil {
		return -1, err
	}

	trans, err := repo.store.Begin()
	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = trans.QueryRow(createQuery, aw.Name, aw.Description, aw.Price, convertRGBAToUint64(aw.Color),
		aw.CreatorId, app.DefaultImage).
		Scan(&aw.ID); err != nil {
		_ = trans.Rollback()
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = repo.setLevel(trans, aw.ID, aw.Price, aw.CreatorId); err != nil {
		_ = trans.Rollback()
		return app.InvalidInt, err
	}

	if err = trans.Commit(); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return aw.ID, nil
}

// GetByID Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetByID(awardsID int64) (*models.Award, error) {
	aw := &models.Award{ID: awardsID}
	var clr uint64
	var childId sql.NullInt64
	if err := repo.store.QueryRow(getByIdQuery, awardsID).
		Scan(&aw.Name, &aw.Description, &aw.Price, &clr, &aw.CreatorId, &aw.Cover, &childId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}
	aw.Color = convertUint64ToRGBA(clr)
	aw.ChildAward = repository.NoAwards
	if childId.Valid {
		aw.ChildAward = childId.Int64
	}
	return aw, nil
}

// CheckAwards Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) CheckAwards(awardsID int64) (bool, error) {
	if err := repo.store.QueryRow(checkAwardsQuery, awardsID).Scan(&awardsID); err != nil {
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
//		repository_postgresql.PriceAlreadyExist
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) Update(aw *models.Award) error {
	creatorId := int64(0)

	if err := repo.store.QueryRow(queryGetCreatorId, aw.ID).Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if err := repo.checkUniq(aw.Name, creatorId, aw.ID, aw.Price); err != nil {
		return err
	}

	trans, err := repo.store.Begin()
	if err != nil {
		return repository.NewDBError(err)
	}

	if _, err = trans.Exec(updateQueryUpdate, aw.Name, aw.Description, aw.Price,
		convertRGBAToUint64(aw.Color), aw.ID); err != nil {
		_ = trans.Rollback()
		return repository.NewDBError(err)
	}

	if err = repo.deleteLevel(trans, aw.ID); err != nil {
		_ = trans.Rollback()
		return err
	}

	if err = repo.setLevel(trans, aw.ID, aw.Price, creatorId); err != nil {
		_ = trans.Rollback()
		return err
	}

	if err = trans.Commit(); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// GetAwards Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) GetAwards(creatorId int64) ([]models.Award, error) {
	var res []models.Award

	rows, err := repo.store.Query(getAwardsQuery, creatorId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		var awards models.Award
		var clr uint64
		var childId sql.NullInt64
		if err = rows.Scan(&awards.ID, &awards.Name, &awards.Description, &awards.Price, &clr,
			&awards.Cover, &childId); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		awards.Color = convertUint64ToRGBA(clr)
		awards.ChildAward = repository.NoAwards
		if childId.Valid {
			awards.ChildAward = childId.Int64
		}
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
	if _, err := repo.store.Exec(deleteQueryUpdate, awardsId); err != nil {
		return repository.NewDBError(err)
	}

	if _, err := repo.store.Exec(deleteQueryDelete, awardsId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// FindByName Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AwardsRepository) FindByName(creatorID int64, awardName string) (bool, error) {
	cnt := 0
	res := repo.store.QueryRow(findByNameQuery, creatorID, awardName).Scan(&cnt)
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
	creatorId := int64(0)

	if err := repo.store.QueryRow(queryGetCreatorId, awardsId).Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	if _, err := repo.store.Exec(updateCoverQuery, cover, awardsId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// setLevel Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) setLevel(trans *sql.Tx, awardsId int64, price int64, creatorId int64) error {
	if _, err := trans.Exec(setLevelQueryInsertParent, awardsId, price, creatorId); err != nil {
		return repository.NewDBError(err)
	}

	if _, err := trans.Exec(setLevelQueryInsertChild, awardsId, price, creatorId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// deleteLevel Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AwardsRepository) deleteLevel(trans *sql.Tx, awardsId int64) error {
	if _, err := trans.Exec(deleteLevelQuery, awardsId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
