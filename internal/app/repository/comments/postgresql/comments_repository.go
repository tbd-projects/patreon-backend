package repository_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"patreon/internal/app/repository/comments"
	postgresql_utilits "patreon/internal/app/utilits/postgresql"

	"github.com/pkg/errors"
)

const (
	checkExistsQuery = "SELECT count(*) from comments where comments_id = $1"

	checkExistsWithPostQuery = "SELECT count(*) from comments where as_creator = $1 and posts_id = $2 and users_id = $3"

	createQuery          = `INSERT INTO comments (body, posts_id, users_id, as_creator) VALUES ($1, $2, $3, $4) RETURNING comments_id`
	createQueryAddToPost = "UPDATE posts SET number_comments = number_comments + 1 where posts_id = $1"

	updateQuery = "UPDATE comments SET body = $1, as_creator = $2 WHERE comments_id = $3"

	getQuery = `SELECT cm.body, cm.as_creator, cm.users_id, cm.posts_id FROM comments AS cm WHERE cm.comments_id = $1`

	getCommentsPostQuery = `
					SELECT cm.comments_id, cm.body, cm.as_creator, cm.users_id, usr.nickname, cm.date, 
					       	CASE WHEN cm.as_creator = TRUE THEN cp.avatar
							ELSE usr.avatar
							END
					FROM comments AS cm
					JOIN users as usr on usr.users_id = cm.users_id
					JOIN creator_profile as cp on cp.creator_id = cm.users_id
					WHERE cm.posts_id = $1
					ORDER BY cm.date DESC LIMIT $2 OFFSET $3;`

	getCommentsUserQuery = `
					SELECT cm.comments_id, cm.body, cm.as_creator, cm.posts_id, ps.title, ps.cover, cm.date
					FROM comments AS cm
					JOIN posts as ps on ps.posts_id = cm.posts_id
					WHERE cm.users_id = $1
					ORDER BY cm.date DESC LIMIT $2 OFFSET $3;`

	deleteQueryDeleteFromPost = "UPDATE posts SET number_comments = number_comments - 1 where posts_id = $1"
	deleteQueryDelete         = "DELETE FROM comments WHERE comments_id = $1 RETURNING posts_id"
)

type CommentsRepository struct {
	store *sqlx.DB
}

var _ = repository_comments.Repository(&CommentsRepository{})

func NewCommentsRepository(st *sqlx.DB) *CommentsRepository {
	return &CommentsRepository{
		store: st,
	}
}

// Create Errors:
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CommentsRepository) Create(cm *models.Comment) (int64, error) {
	trans, err := repo.store.Begin()
	if err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	if err = trans.QueryRow(createQuery, cm.Body, cm.PostId, cm.AuthorId, cm.AsCreator).
		Scan(&cm.ID); err != nil {
		_ = trans.Rollback()
		return app.InvalidInt,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for post %d", cm.PostId)))
	}

	if _, err = trans.Exec(createQueryAddToPost, cm.PostId); err != nil {
		_ = trans.Rollback()
		return app.InvalidInt,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments %d for post %d",
				cm.ID, cm.PostId)))
	}

	if err = trans.Commit(); err != nil {
		return app.InvalidInt, repository.NewDBError(err)
	}

	return cm.ID, nil
}

// Update Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CommentsRepository) Update(cm *models.Comment) error {
	if err := repo.CheckExists(cm.ID); !errors.Is(err, CommentAlreadyExist) {
		return err
	}

	if _, err := repo.store.Exec(updateQuery, cm.Body, cm.AsCreator, cm.ID); err != nil {
		return repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments %d with body %s",
			cm.ID, cm.Body)))
	}

	return nil
}

// GetUserComments Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CommentsRepository) GetUserComments(userId int64, pag *models.Pagination) ([]models.UserComment, error) {
	limit, offset, er := postgresql_utilits.AddPagination("comments", pag, repo.store)
	if er != nil {
		return nil, er
	}

	rows, err := repo.store.Query(getCommentsUserQuery, userId, limit, offset)
	if err != nil {
		return nil,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for user %d", userId)))
	}

	var comments []models.UserComment
	for rows.Next() {
		comment := models.UserComment{Comment: models.Comment{AuthorId: userId}}
		if err = rows.Scan(&comment.ID, &comment.Body, &comment.AsCreator, &comment.PostId,
			&comment.PostName, &comment.PostCover, &comment.Date); err != nil {
			_ = rows.Close()
			return nil,
				repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for user %d", userId)))
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for user %d", userId)))
	}
	return comments, nil
}

// GetPostComments Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CommentsRepository) GetPostComments(postId int64, pag *models.Pagination) ([]models.PostComment, error) {
	limit, offset, er := postgresql_utilits.AddPagination("comments", pag, repo.store)
	if er != nil {
		return nil, er
	}

	rows, err := repo.store.Query(getCommentsPostQuery, postId, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for post %d", postId)))
	}

	var comments []models.PostComment
	for rows.Next() {
		comment := models.PostComment{Comment: models.Comment{PostId: postId}}
		if err = rows.Scan(&comment.ID, &comment.Body, &comment.AsCreator, &comment.AuthorId,
			&comment.AuthorNickname, &comment.Date, &comment.AuthorAvatar); err != nil {
			_ = rows.Close()
			return nil,
				repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for post %d", postId)))
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil,
			repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try create comments for post %d", postId)))
	}
	return comments, nil
}

// CheckExists Errors:
//		repository_postgresql.CommentAlreadyExist
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CommentsRepository) CheckExists(commentId int64) error {
	cnt := int64(0)
	if err := repo.store.Get(&cnt, checkExistsQuery, commentId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(errors.Wrap(err,
			fmt.Sprintf("try checkExists comments %d", commentId)))
	}

	if cnt != 0 {
		return CommentAlreadyExist
	}

	return repository.NotFound
}

// Get Errors:
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CommentsRepository) Get(commentsId int64) (*models.Comment, error) {
	cm := &models.Comment{ID: commentsId}
	if err := repo.store.QueryRowx(getQuery, commentsId).
		Scan(&cm.Body, &cm.AsCreator, &cm.AuthorId, &cm.PostId); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, repository.NewDBError(errors.Wrap(err,
			fmt.Sprintf("try checkExists comments with id %d", commentsId)))
	}

	return cm, nil
}

// CheckExists Errors:
//		repository_postgresql.CommentAlreadyExist
//		repository.NotFound
// 		app.GeneralError with Errors
// 			repository.DefaultErrDB
func (repo *CommentsRepository) checkExistsWithPost(authorId int64, postId int64, asCreator bool) error {
	cnt := int64(0)
	if err := repo.store.Get(&cnt, checkExistsWithPostQuery, asCreator, postId, authorId); err != nil {
		if err == sql.ErrNoRows {
			return repository.NotFound
		}
		return repository.NewDBError(errors.Wrap(err,
			fmt.Sprintf("try checkExists comments with post %d and user %d", authorId, postId)))
	}

	if cnt != 0 {
		return CommentAlreadyExist
	}

	return repository.NotFound
}

// Delete Errors:
// 		app.GeneralError with Errors:
// 			repository.DefaultErrDB
func (repo *CommentsRepository) Delete(commentId int64) error {
	tx, err := repo.store.Beginx()
	if err != nil {
		return repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try delete comments %d", commentId)))
	}

	postId := int64(0)
	if err = tx.Get(&postId, deleteQueryDelete, commentId); err != nil {
		_ = tx.Rollback()
		return repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try delete comments %d", commentId)))
	}

	if _, err = tx.Exec(deleteQueryDeleteFromPost, postId); err != nil {
		_ = tx.Rollback()
		return repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try delete comments %d from post %d",
			commentId, postId)))
	}

	if err = tx.Commit(); err != nil {
		return repository.NewDBError(errors.Wrap(err, fmt.Sprintf("try delete comments %d", commentId)))
	}

	return nil
}
