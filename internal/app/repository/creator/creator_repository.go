package repository_creator

import (
	"database/sql"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
)

type CreatorRepository struct {
	store *sql.DB
}

func NewCreatorRepository(st *sql.DB) *CreatorRepository {
	return &CreatorRepository{
		store: st,
	}
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CreatorRepository) Create(cr *models.Creator) (int64, error) {
	if err := repo.store.QueryRow("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING creator_id", cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).Scan(&cr.ID); err != nil {
		return -1, repository.NewDBError(err)
	}
	return cr.ID, nil
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreators() ([]models.Creator, error) {
	count := 0

	if err := repo.store.QueryRow("SELECT count(*) from creator_profile").Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.Query(
		"SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
			"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			return nil, repository.NewDBError(err)
		}
		res[i] = creator
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

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreator(creatorId int64) (*models.Creator, error) {
	creator := &models.Creator{}

	if err := repo.store.QueryRow("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname "+
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1", creatorId).
		Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return creator, nil
}
