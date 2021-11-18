package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"strings"
	"time"
)

const UnusedAttach = -1 //value stored in level field

const (
	deleteUnusedQuery = `DELETE FROM posts_data WHERE level = $1 and post_id = $2`

	makeUnusedAttachQuery = `UPDATE posts_data SET level = $1 WHERE post_id = $2`

	getDataTypeAndIdQuery = `SELECT posts_type_id, type FROM posts_type`

	updateAttachQuery      = `UPDATE posts_data SET data = $1, level = $2 WHERE data_id = $3 RETURNING data_id`
	updateAttachFilesQuery = `UPDATE posts_data SET level = $1 WHERE data_id = $2 RETURNING data_id`

	createAttachesQuery    = `INSERT INTO posts_data (post_id, type, data, level) VALUES`
	createAttachesQueryEnd = `RETURNING data_id`
)

// getAttachTypeAndId Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AttachesRepository) getAttachTypeAndId() (map[models.DataType]int64, error) {
	if repo.dataTypes != nil && repo.lastUpdate.Add(reloadDataType).Before(time.Now()) {
		return repo.dataTypes, nil
	}
	repo.lastUpdate = time.Now()
	tmpDataTypes := map[models.DataType]int64{}

	rows, err := repo.store.Queryx(getDataTypeAndIdQuery)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	for rows.Next() {
		dataType := models.DataType("")
		attachId := int64(3)
		if err = rows.Scan(&attachId, &dataType); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		tmpDataTypes[dataType] = attachId
	}

	repo.dataTypes = tmpDataTypes
	return repo.dataTypes, nil
}

// createAttaches Errors:
//		UnknownDataFormat
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *AttachesRepository) createAttaches(trans *sqlx.Tx, postId int64,
	newAttachs []models.Attach) ([]models.Attach, error) {
	dataTypes, err := repo.getAttachTypeAndId()
	if err != nil {
		return nil, err
	}

	var argsString []string
	var args []interface{}
	bdIndex := 1
	for _, attach := range newAttachs {
		argsString = append(argsString, "(?, ?, ?, ?)")

		args = append(args, postId)
		if _, ok := dataTypes[attach.Type]; !ok {
			return nil, errors.Wrap(UnknownDataFormat, fmt.Sprintf("err with type %s", attach.Type))
		}
		args = append(args, dataTypes[attach.Type])
		args = append(args, attach.Value)
		args = append(args, attach.Level)

		bdIndex += 4
	}

	query := fmt.Sprintf("%s %s %s", createAttachesQuery,
		strings.Join(argsString, ", "), createAttachesQueryEnd)
	query = repo.store.Rebind(query)

	var res []int64
	if err = trans.Select(&res, query, args...); err != nil {
		return nil, repository.NewDBError(err)
	}

	for index, id := range res {
		newAttachs[index].Id = id
	}
	return newAttachs, nil
}

// ApplyChangeAttaches Errors:
//		UnknownDataFormat
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) ApplyChangeAttaches(postId int64,
	newAttaches []models.Attach, updatedAttaches []models.Attach) ([]int64, error) {
	res := make([]int64, len(newAttaches)+len(updatedAttaches))

	trans, err := repo.store.Beginx()
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	if err = repo.markUnusedAttach(trans, postId); err != nil {
		_ = trans.Rollback()
		return nil, err
	}

	if len(newAttaches) != 0 {
		if newAttaches, err = repo.createAttaches(trans, postId, newAttaches); err != nil {
			_ = trans.Rollback()
			return nil, err
		}
	}

	for _, attach := range updatedAttaches {
		if err = repo.updateAttach(trans, &attach); err != nil {
			break
		}
		res[attach.Level-1] = attach.Id
	}

	if err != nil {
		_ = trans.Rollback()
		return nil, err
	}

	if err = repo.deleteUnused(trans, postId); err != nil {
		_ = trans.Rollback()
		return nil, err
	}

	if err = trans.Commit(); err != nil {
		return nil, err
	}

	for _, attach := range newAttaches {
		res[attach.Level-1] = attach.Id
	}

	return res, nil
}

// updateAttach Errors:
//		UnknownDataFormat
//		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) updateAttach(trans *sqlx.Tx, attach *models.Attach) error {
	dataTypes, err := repo.getAttachTypeAndId()
	if err != nil {
		return err
	}

	if _, ok := dataTypes[attach.Type]; !ok {
		return UnknownDataFormat
	}

	if attach.Type == models.Text {
		err = trans.QueryRow(updateAttachQuery, attach.Value, attach.Level, attach.Id).Scan(&attach.Id)
	} else {
		err = trans.QueryRow(updateAttachFilesQuery, attach.Level, attach.Id).Scan(&attach.Id)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// DeleteUnused Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) deleteUnused(trans *sqlx.Tx, postId int64) error {
	_, err := trans.Exec(deleteUnusedQuery, UnusedAttach, postId)
	if err != nil {
		return repository.NewDBError(err)
	}

	return nil
}

// markUnusedAttach Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *AttachesRepository) markUnusedAttach(trans *sqlx.Tx, postId int64) error {
	_, err := trans.Exec(makeUnusedAttachQuery, UnusedAttach, postId)
	if err != nil {
		return repository.NewDBError(err)
	}

	return nil
}
