package repository_postgresql

import (
	"github.com/jmoiron/sqlx"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

const (
	queryCategoryGet = `SELECT DISTINCT name FROM creator_category ORDER BY name`
	queryTypeDataGet = `SELECT DISTINCT type FROM posts_type ORDER BY type`
)

type InfoRepository struct {
	store *sqlx.DB
}

func NewInfoRepository(st *sqlx.DB) *InfoRepository {
	return &InfoRepository{
		store: st,
	}
}

// Get Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *InfoRepository) Get() (*models.Info, error) {
	rowCaregory, err := repo.store.Query(queryCategoryGet)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	var categories []string
	for rowCaregory.Next() {
		category := ""
		if err = rowCaregory.Scan(&category); err != nil {
			_ = rowCaregory.Close()
			return nil, repository.NewDBError(err)
		}
		categories = append(categories, category)
	}

	if err = rowCaregory.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	rowType, err := repo.store.Query(queryTypeDataGet)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	var dataTypes []string

	for rowType.Next() {
		dataType := ""
		if err = rowType.Scan(&dataType); err != nil {
			_ = rowType.Close()
			return nil, repository.NewDBError(err)
		}
		dataTypes = append(dataTypes, dataType)
	}

	if err = rowType.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return &models.Info{Category: categories, TypePostData: dataTypes}, nil
}
