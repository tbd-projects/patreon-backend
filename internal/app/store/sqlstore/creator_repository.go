package sqlstore

import (
	"database/sql"
	"patreon/internal/app/store"
	"patreon/internal/models"

	log "github.com/sirupsen/logrus"
)

type CreatorRepository struct {
	store          *Store
	UserRepository UserRepository
}

func NewCreatorRepository(st *Store) *CreatorRepository {
	return &CreatorRepository{
		store:          st,
		UserRepository: UserRepository{st},
	}
}

func (repo *CreatorRepository) Create(cr *models.Creator) error {
	if err := repo.store.db.QueryRow("INSERT INTO creator_profile (creator_id, category, "+
		"description, avatar, cover) VALUES ($1, $2, $3, $4, $5)"+
		"RETURNING creator_id", cr.ID, cr.Category, cr.Description, cr.Avatar, cr.Cover).Scan(&cr.ID); err != nil {
		return err
	}
	return nil
}

func (repo *CreatorRepository) GetCreators() ([]models.Creator, error) {
	count := 0

	if err := repo.store.db.QueryRow("SELECT count(*) from creator_profile").Scan(&count); err != nil {
		return nil, err
	}
	res := make([]models.Creator, count)

	rows, err := repo.store.db.Query(
		"SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname " +
			"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id")
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	i := 0
	for rows.Next() {
		var creator models.Creator
		if err = rows.Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
			return nil, err
		}
		res[i] = creator
		i++

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return res, err
}

func (repo *CreatorRepository) GetCreator(creatorId int64) (*models.Creator, error) {
	creator := &models.Creator{}

	if err := repo.store.db.QueryRow("SELECT creator_id, category, description, creator_profile.avatar, cover, usr.nickname "+
		"from creator_profile join users as usr on usr.user_id = creator_profile.creator_id where creator_id=$1", creatorId).
		Scan(&creator.ID, &creator.Category, &creator.Description, &creator.Avatar,
			&creator.Cover, &creator.Nickname); err != nil {
		return nil, store.NotFound
	}

	return creator, nil
}
