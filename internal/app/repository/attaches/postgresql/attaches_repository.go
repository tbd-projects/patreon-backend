package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_attaches "patreon/internal/app/repository/attaches"
	"time"
)

const reloadDataType = 24 * time.Hour

const (
	getAndCheckDataTypeIdQuery = `SELECT posts_type_id FROM posts_type WHERE type = $1`

	getDataTypeQuery = `SELECT type FROM posts_type WHERE posts_type_id = $1`

	createQuery = `INSERT INTO posts_data (type, data, post_id) VALUES ($1, $2, $3) 
		RETURNING data_id`

	getQuery = `SELECT post_id, data, type FROM posts_data WHERE data_id = $1`

	existsAttachQuery = `SELECT post_id FROM posts_data WHERE data_id in (?)`

	getAttachesQuery = `SELECT data_id, pst.type, data FROM posts_data JOIN posts_type AS pst 
    			ON (pst.posts_type_id = posts_data.type) WHERE post_id = $1 ORDER BY level`

	updateQuery = `UPDATE posts_data SET type = $1, data = $2 WHERE data_id = $3 RETURNING data_id`

	deleteQuery = `DELETE FROM posts_data WHERE data_id = $1`
)

type AttachesRepository struct {
	store      *sqlx.DB
	dataTypes  map[models.DataType]int64
	lastUpdate time.Time
}

var _ = repository_attaches.Repository(&AttachesRepository{})

func NewAttachesRepository(st *sqlx.DB) *AttachesRepository {
	return &AttachesRepository{
		store:      st,
		lastUpdate: time.Now(),
		dataTypes:  nil,
	}
}

// getAndCheckDataTypeId Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AttachesRepository) getAndCheckAttachTypeId(dataType models.DataType) (int64, error) {
	var dataTypeId int64
	if err := repo.store.QueryRow(getAndCheckDataTypeIdQuery, dataType).
		Scan(&dataTypeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, errors.Wrap(UnknownDataFormat, fmt.Sprintf("err with type %s", dataType))
		}
		return app.InvalidInt, repository.NewDBError(err)
	}
	return dataTypeId, nil
}

// getDataType Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AttachesRepository) getAttachType(dataTypeId int64) (models.DataType, error) {
	var dataType models.DataType
	if err := repo.store.QueryRow(getDataTypeQuery, dataTypeId).
		Scan(&dataType); err != nil {
		return "", repository.NewDBError(err)
	}
	return dataType, nil
}

// Create Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AttachesRepository) Create(postData *models.AttachWithoutLevel) (int64, error) {
	type_id, err := repo.getAndCheckAttachTypeId(postData.Type)
	if err != nil {
		return app.InvalidInt, err
	}

	if err = repo.store.QueryRow(createQuery, type_id, postData.Value, postData.PostId).
		Scan(&postData.ID); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return postData.ID, nil
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) Get(attachId int64) (*models.AttachWithoutLevel, error) {
	data := &models.AttachWithoutLevel{ID: attachId}
	var typeId int64
	if err := repo.store.QueryRow(getQuery, attachId).Scan(&data.PostId, &data.Value,
		&typeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}
	var err error
	if data.Type, err = repo.getAttachType(typeId); err != nil {
		return nil, err
	}

	return data, nil
}

// ExistsAttach Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) ExistsAttach(attachId ...int64) (bool, error) {
	if attachId == nil {
		return true, nil
	}
	query, args, err := sqlx.In(existsAttachQuery, attachId)
	if err != nil {
		return false, repository.NewDBError(err)
	}
	query = repo.store.Rebind(query)
	var selectedData []int64
	if err = repo.store.Select(&selectedData, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.NotFound
		}
		return false, repository.NewDBError(err)
	}

	return len(selectedData) == len(attachId), nil
}

// GetAttaches Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) GetAttaches(postsId int64) ([]models.AttachWithoutLevel, error) {
	var res []models.AttachWithoutLevel

	rows, err := repo.store.Query(getAttachesQuery, postsId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		var data models.AttachWithoutLevel
		if err = rows.Scan(&data.ID, &data.Type, &data.Value); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}

		res = append(res, data)
		data.PostId = postsId
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// Update Errors:
//		UnknownDataFormat
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) Update(postData *models.AttachWithoutLevel) error {
	type_id, err := repo.getAndCheckAttachTypeId(postData.Type)
	if err != nil {
		return err
	}

	if err = repo.store.QueryRow(updateQuery, type_id, postData.Value, postData.ID).
		Scan(&postData.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}
	return nil
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) Delete(attachId int64) error {
	_, err := repo.store.Exec(deleteQuery, attachId)
	if err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
