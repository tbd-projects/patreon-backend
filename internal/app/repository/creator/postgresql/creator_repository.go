package repository_postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"patreon/internal/app"
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
//		IncorrectCategory
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CreatorRepository) Create(cr *models.Creator) (int64, error) {
	queryCategory := `SELECT category_id FROM creator_category WHERE name = $1`

	query := `INSERT INTO creator_profile (creator_id, category,
		description, avatar, cover) VALUES ($1, $2, $3, $4, $5)
		RETURNING creator_id
	`

	category := int64(0)
	if err := repo.store.QueryRow(queryCategory, cr.Category).Scan(&category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, IncorrectCategory
		}
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err := repo.store.QueryRow(query, cr.ID, category, cr.Description, cr.Avatar, cr.Cover).Scan(&cr.ID); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}
	return cr.ID, nil
}

// GetCreators Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreators() ([]models.Creator, error) {
	queryCount := `SELECT count(*) from creator_profile`
	queryCreator := `SELECT creator_id, cc.name, description, creator_profile.avatar, cover, usr.nickname 
					FROM creator_profile JOIN users AS usr ON usr.users_id = creator_profile.creator_id
					JOIN creator_category As cc ON creator_profile.category = cc.category_id`
	count := 0

	if err := repo.store.QueryRow(queryCount).Scan(&count); err != nil {
		return nil, repository.NewDBError(err)
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.Query(queryCreator)
	if err != nil {
		return nil, repository.NewDBError(err)
	}

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			_ = rows.Close()
			return nil, repository.NewDBError(err)
		}
		res[i] = creator
		i++
	}

	if err = rows.Err(); err != nil {
		return nil, repository.NewDBError(err)
	}

	return res, nil
}

// GetCreator Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) GetCreator(creatorId int64) (*models.Creator, error) {
	query := `SELECT creator_id, cc.name, description, creator_profile.avatar, cover, usr.nickname 
			FROM creator_profile JOIN users AS usr ON usr.users_id = creator_profile.creator_id 
			JOIN creator_category As cc ON creator_profile.category = cc.category_id 
			where creator_id=$1`
	creator := &models.Creator{}

	if err := repo.store.QueryRow(query, creatorId).
		Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(err)
	}

	return creator, nil
}

// ExistsCreator Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) ExistsCreator(creatorId int64) (bool, error) {
	query := `SELECT creator_id from creator_profile where creator_id=$1`
	creator := &models.Creator{}

	if err := repo.store.QueryRow(query, creatorId).Scan(&creator.ID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, repository.NewDBError(err)
	}

	return true, nil
}

// UpdateAvatar Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) UpdateAvatar(creatorId int64, avatar string) error {
	query := `UPDATE creator_profile SET avatar = $1 WHERE creator_id = $2 RETURNING creator_id`

	if err := repo.store.QueryRow(query, avatar, creatorId).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}

// UpdateCover Errors:
// 		repository.NotFound
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CreatorRepository) UpdateCover(creatorId int64, cover string) error {
	query := `UPDATE creator_profile SET cover = $1 WHERE creator_id = $2 RETURNING creator_id`

	if err := repo.store.QueryRow(query, cover, creatorId).
		Scan(&creatorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(err)
	}

	return nil
}
