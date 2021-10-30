package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	repository_posts_data "patreon/internal/app/repository/posts_data"
)

type PostsDataRepository struct {
	store *sql.DB
}

var _ = repository_posts_data.Repository(&PostsDataRepository{})

func NewPostsDataRepository(st *sql.DB) *PostsDataRepository {
	return &PostsDataRepository{
		store: st,
	}
}

// getAndCheckDataTypeId Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *PostsDataRepository) getAndCheckDataTypeId(dataType models.DataType) (int64, error) {
	query := `SELECT posts_type_id FROM posts_type WHERE type = $1`
	var dataTypeId int64
	if err := repo.store.QueryRow(query, dataType).
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
func (repo *PostsDataRepository) getDataType(dataTypeId int64) (models.DataType, error) {
	query := `SELECT type FROM posts_type WHERE posts_type_id = $1`
	var dataType models.DataType
	if err := repo.store.QueryRow(query, dataTypeId).
		Scan(&dataType); err != nil {
		return "", repository.NewDBError(err)
	}
	return dataType, nil
}

// Create Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *PostsDataRepository) Create(postData *models.PostData) (int64, error) {
	query := `INSERT INTO posts_data (type, data, post_id) VALUES ($1, $2, $3) 
		RETURNING data_id`

	type_id, err := repo.getAndCheckDataTypeId(postData.Type)
	if err != nil {
		return app.InvalidInt, err
	}

	if err = repo.store.QueryRow(query, type_id, postData.Data, postData.PostId).
		Scan(&postData.ID); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return postData.ID, nil
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsDataRepository) Get(dataID int64) (*models.PostData, error) {
	query := `SELECT post_id, data, type FROM posts_data WHERE data_id = $1`

	data := &models.PostData{ID: dataID}
	var typeId int64
	if err := repo.store.QueryRow(query, dataID).Scan(&data.PostId, &data.Type,
		&typeId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}
	var err error
	if data.Type, err = repo.getDataType(typeId); err != nil {
		return nil, err
	}

	return data, nil
}

// GetData Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsDataRepository) GetData(postsId int64) ([]models.PostData, error) {
	query := `SELECT data_id, pst.type, data FROM posts_data JOIN posts_type AS pst 
    			ON (pst.posts_type_id = posts_data.type) WHERE post_id = $1`

	var res []models.PostData

	rows, err := repo.store.Query(query, postsId)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var data models.PostData
		if err = rows.Scan(&data.ID, &data.Type, &data.Data); err != nil {
			return nil, repository.NewDBError(err)
		}
		res[i] = data
		data.PostId = postsId
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

// Update Errors:
//		UnknownDataFormat
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *PostsDataRepository) Update(postData *models.PostData) error {
	query := `UPDATE posts_data SET type = $1, data = $2, post_id = $3 WHERE data_id = $4`

	type_id, err := repo.getAndCheckDataTypeId(postData.Type)
	if err != nil {
		return err
	}

	if err = repo.store.QueryRow(query, type_id, postData.Data, postData.PostId).
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
func (repo *PostsDataRepository) Delete(dataId int64) error {
	query := `DELETE FROM posts_data WHERE data_id = $q`

	if _, err := repo.store.Query(query, dataId); err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
