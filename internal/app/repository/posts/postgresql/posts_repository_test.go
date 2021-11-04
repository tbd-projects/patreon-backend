package repository_postgresql

import (
	"database/sql"
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	putilits "patreon/internal/app/utilits/postgresql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"
)

type SuitePostsRepository struct {
	models.Suite
	repo *PostsRepository
}

func (s *SuitePostsRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewPostsRepository(s.DB)
}

func (s *SuitePostsRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuitePostsRepository) TestPostsRepository_Create() {
	query := `INSERT INTO posts (title, description,
		type_awards, creator_id, cover) VALUES ($1, $2, $3, $4, $5) 
		RETURNING posts_id`
	post := &models.CreatePost{ID: 2, Title: "sad", Description: "asdasd", Awards: 1, CreatorId: 2}

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, post.Awards, post.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"posts_id"}).AddRow(post.ID))
	id, err := s.repo.Create(post)
	assert.Equal(s.T(), post.ID, id)
	assert.NoError(s.T(), err)

	post.Awards = repository.NoAwards
	var awardsId sql.NullInt64
	awardsId.Int64 = repository.NoAwards
	awardsId.Valid = false

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, awardsId, post.CreatorId, app.DefaultImage).
		WillReturnRows(sqlmock.NewRows([]string{"posts_id"}).AddRow(post.ID))
	id, err = s.repo.Create(post)
	assert.Equal(s.T(), post.ID, id)
	assert.NoError(s.T(), err)

	post.Awards = 1
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, post.Awards, post.CreatorId, app.DefaultImage).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.Create(post)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func (s *SuitePostsRepository) TestPostsRepository_GetPostCreator() {
	query := `SELECT creator_id FROM posts WHERE posts_id = $1`
	creatorId := int64(3)
	postId := int64(2)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{"creator_id"}).AddRow(creatorId))
	id, err := s.repo.GetPostCreator(postId)
	assert.Equal(s.T(), id, creatorId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.GetPostCreator(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetPostCreator(postId)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsRepository) TestPostsRepository_GetPost() {
	query := `
			SELECT title, description, likes, posts.date, cover, type_awards, creator_id, lk.likes_id IS NOT NULL, views FROM posts
				LEFT OUTER JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
				WHERE posts.posts_id = $2;`
	queryPost := `UPDATE posts SET views = views + 1 WHERE posts_id = $1`

	post := &models.Post{ID: 2, Title: "sad", Description: "asdasd", Awards: 1, CreatorId: 2}
	userId := int64(5)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				post.Awards, post.CreatorId, post.AddLike, post.Views))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryPost)).
		WithArgs(post.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	res, err := s.repo.GetPost(post.ID, userId, true)
	assert.Equal(s.T(), res, post)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				post.Awards, post.CreatorId, post.AddLike, post.Views))
	res, err = s.repo.GetPost(post.ID, userId, false)
	assert.Equal(s.T(), res, post)
	assert.NoError(s.T(), err)

	var awardsId sql.NullInt64
	awardsId.Int64 = repository.NoAwards
	awardsId.Valid = false

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				awardsId, post.CreatorId, post.AddLike, post.Views))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryPost)).
		WithArgs(post.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	res, err = s.repo.GetPost(post.ID, userId, true)
	post.Awards = repository.NoAwards
	assert.Equal(s.T(), res, post)
	post.Awards = 1
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				post.Awards, post.CreatorId, post.AddLike, post.Views))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryPost)).
		WithArgs(post.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow().CloseError(models.BDError))
	_, err = s.repo.GetPost(post.ID, userId, true)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				post.Awards, post.CreatorId, post.AddLike, post.Views))
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryPost)).
		WithArgs(post.ID).
		WillReturnError(models.BDError)
	_, err = s.repo.GetPost(post.ID, userId, true)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnError(models.BDError)
	_, err = s.repo.GetPost(post.ID, userId, true)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.ID).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetPost(post.ID, userId, true)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsRepository) TestPostsRepository_GetPosts() {
	queryStat := "SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = $1"
	tableName := "posts"

	pag := &models.Pagination{Limit: 10, Offset: 20}
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	limit, offset, err := putilits.AddPagination(tableName, pag, s.DB)
	assert.NoError(s.T(), err)
	query := `
			SELECT posts_id, title, description, likes, type_awards, posts.date, cover, lk.likes_id IS NOT NULL, views
			FROM posts
			LEFT JOIN likes AS lk ON (lk.post_id = posts.posts_id and lk.users_id = $1)
			WHERE creator_id = $2 ORDER BY posts.date ` + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	post := models.Post{ID: 2, Title: "sad", Description: "asdasd", Awards: 1, CreatorId: 2}
	userId := int64(5)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.CreatorId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "title", "description", "likes",
			"type_awards", "posts.date", "cover", "have_like", "views"}).
			AddRow(post.ID, post.Title, post.Description, post.Likes, post.Awards, post.Date, post.Cover,
				post.AddLike, post.Views))
	res, err := s.repo.GetPosts(post.CreatorId, userId, pag)
	assert.Equal(s.T(), res[0], post)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnError(repository.DefaultErrDB)
	_, err = s.repo.GetPosts(post.CreatorId, userId, pag)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	var awardsId sql.NullInt64
	awardsId.Int64 = repository.NoAwards
	awardsId.Valid = false

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.CreatorId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "title", "description", "likes",
			"type_awards", "posts.date", "cover", "have_like", "views"}).
			AddRow(post.ID, post.Title, post.Description, post.Likes, awardsId, post.Date, post.Cover,
				post.AddLike, post.Views))
	res, err = s.repo.GetPosts(post.CreatorId, userId, pag)
	post.Awards = repository.NoAwards
	assert.Equal(s.T(), res[0], post)
	post.Awards = 1
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.CreatorId).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "likes",
			"posts.date", "cover", "type_awards", "creator_id", "have_like", "views"}).
			AddRow(post.Title, post.Description, post.Likes, post.Date, post.Cover,
				post.Awards, post.CreatorId, post.Description, post.Views))
	_, err = s.repo.GetPosts(post.CreatorId, userId, pag)
	assert.Error(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.CreatorId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "title", "description", "likes",
			"type_awards", "posts.date", "cover", "have_like", "views"}).
			AddRow(post.ID, post.Title, post.Description, post.Likes, post.Awards, post.Date, post.Cover,
				post.AddLike, post.Views).RowError(0, models.BDError))
	_, err = s.repo.GetPosts(post.CreatorId, userId, pag)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryStat)).
		WithArgs(tableName).
		WillReturnRows(sqlmock.NewRows([]string{"n_live_tup"}).AddRow(int64(5000)))
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId, post.CreatorId).
		WillReturnError(models.BDError)
	_, err = s.repo.GetPosts(post.CreatorId, userId, pag)
	assert.Error(s.T(), err, repository.NewDBError(models.BDError))
}

func (s *SuitePostsRepository) TestPostsRepository_Update() {
	query := `UPDATE posts SET title = $1, description = $2, type_awards = $3 WHERE posts_id = $4 RETURNING posts_id`
	post := &models.UpdatePost{ID: 2, Title: "sad", Description: "asdasd", Awards: 1}

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, post.Awards, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"posts_id"}).AddRow(post.ID))
	err := s.repo.UpdatePost(post)
	assert.NoError(s.T(), err)

	post.Awards = repository.NoAwards
	var awardsId sql.NullInt64
	awardsId.Int64 = repository.NoAwards
	awardsId.Valid = false

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, awardsId, post.ID).
		WillReturnRows(sqlmock.NewRows([]string{"posts_id"}).AddRow(post.ID))
	err = s.repo.UpdatePost(post)
	assert.NoError(s.T(), err)

	post.Awards = 1
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, post.Awards, post.ID).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.UpdatePost(post)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(post.Title, post.Description, post.Awards, post.ID).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.UpdatePost(post)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsRepository) TestPostsRepository_UpdateCover() {
	query := `UPDATE posts SET cover = $1 WHERE posts_id = $2 RETURNING posts_id`

	postId := int64(2)
	cover := "sdadsasd"
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, postId).
		WillReturnRows(sqlmock.NewRows([]string{"posts_id"}).AddRow(postId))
	err := s.repo.UpdateCoverPost(postId, cover)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, postId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.UpdateCoverPost(postId, cover)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cover, postId).
		WillReturnError(sql.ErrNoRows)
	err = s.repo.UpdateCoverPost(postId, cover)
	assert.Error(s.T(), err, repository.NotFound)
}

func (s *SuitePostsRepository) TestPostsRepository_Delete() {
	query := `DELETE FROM posts WHERE posts_id = $q`

	postId := int64(2)
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{}))
	err := s.repo.Delete(postId)
	assert.NoError(s.T(), err)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnError(repository.DefaultErrDB)
	err = s.repo.Delete(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(postId).
		WillReturnRows(sqlmock.NewRows([]string{}).CloseError(repository.DefaultErrDB))
	err = s.repo.Delete(postId)
	assert.Error(s.T(), err, repository.NewDBError(repository.DefaultErrDB))
}

func TestPostsRepository(t *testing.T) {
	suite.Run(t, new(SuitePostsRepository))
}
