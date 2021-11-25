package repository_postgresql

import (
	"database/sql"
	"fmt"
	"patreon/internal/app"
	"patreon/internal/app/models"
	"patreon/internal/app/repository"
	"regexp"
	"testing"

	"github.com/lib/pq"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"

	"github.com/zhashkevych/go-sqlxmock"

	"github.com/stretchr/testify/require"
)

type SuiteLikesRepository struct {
	models.Suite
	repo *LikesRepository
}
type LikesWithFieldMatcher struct{ like *models.Like }

func newLikesWithFieldMatcher(like *models.Like) gomock.Matcher {
	return &LikesWithFieldMatcher{like}
}

func (match *LikesWithFieldMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case *models.Like:
		return x.(*models.Like).ID == match.like.ID && x.(*models.Like).UserId == match.like.UserId &&
			x.(*models.Like).Value == match.like.Value && x.(*models.Like).PostId == match.like.PostId
	default:
		return false
	}
}

func (match *LikesWithFieldMatcher) String() string {
	return fmt.Sprintf("Like: %v", match.like)
}
func (s *SuiteLikesRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewLikesRepository(s.DB)
}

func (s *SuiteLikesRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}
func (s *SuiteLikesRepository) TestLikesRepositoryGet_Correct() {
	query := `SELECT post_id, value, likes_id FROM likes WHERE users_id = $1`
	like := TestLike(s.T())

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value", "likes_id"}).
			AddRow(like.PostId, like.Value, like.ID))
	res, err := s.repo.Get(like.UserId)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), like, res)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGet_UserNotFound() {
	query := `SELECT post_id, value, likes_id FROM likes WHERE users_id = $1`
	like := TestLike(s.T())
	expErr := repository.NotFound

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId).
		WillReturnError(sql.ErrNoRows)

	res, err := s.repo.Get(like.UserId)
	assert.Equal(s.T(), err, expErr)
	assert.Nil(s.T(), res)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGet_DBError() {
	query := `SELECT post_id, value, likes_id FROM likes WHERE users_id = $1`
	like := TestLike(s.T())
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId).
		WillReturnError(sqlErr)

	res, err := s.repo.Get(like.UserId)
	assert.Equal(s.T(), err, expErr)
	assert.Nil(s.T(), res)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGetLikeId_DBError() {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	like := TestLike(s.T())

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	resExp := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId, like.PostId).
		WillReturnError(sqlErr)
	res, err := s.repo.GetLikeId(like.UserId, like.PostId)

	assert.Equal(s.T(), resExp, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGetLikeId_SqlNoRows() {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	like := TestLike(s.T())

	sqlErr := sql.ErrNoRows
	expErr := repository.NotFound
	resExp := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId, like.PostId).
		WillReturnError(sqlErr)
	res, err := s.repo.GetLikeId(like.UserId, like.PostId)

	assert.Equal(s.T(), resExp, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGetLikeId_NotFoundUserLikeOnPost() {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	like := TestLike(s.T())

	sqlErr := repository.NotFound
	expErr := repository.NewDBError(sqlErr)
	resExp := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId, like.PostId).
		WillReturnError(sqlErr)
	res, err := s.repo.GetLikeId(like.UserId, like.PostId)

	assert.Equal(s.T(), resExp, res)
	assert.Equal(s.T(), expErr, err)
}

func (s *SuiteLikesRepository) TestLikesRepositoryGetLikeId_InternalError() {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	like := TestLike(s.T())

	sqlErr := repository.NotFound
	expErr := repository.NewDBError(sqlErr)
	resExp := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId, like.PostId).
		WillReturnError(sqlErr)
	res, err := s.repo.GetLikeId(like.UserId, like.PostId)

	assert.Equal(s.T(), resExp, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryGetLikeId_Correct() {
	query := `SELECT likes_id FROM likes WHERE users_id = $1 AND post_id = $2`
	like := TestLike(s.T())

	resExp := like.ID

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(like.UserId, like.PostId).
		WillReturnRows(sqlmock.NewRows([]string{"likes_id"}).
			AddRow(like.ID))

	res, err := s.repo.GetLikeId(like.UserId, like.PostId)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), resExp, res)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_BeginTransactionError() {

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(pq.ErrNotSupported)
	expRes := int64(app.InvalidInt)

	like := TestLike(s.T())

	s.Mock.ExpectBegin().WillReturnError(sqlErr)

	res, err := s.repo.Add(like)

	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_InsertQueryError() {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(pq.ErrNotSupported)
	expRes := int64(app.InvalidInt)

	like := TestLike(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
		WithArgs(like.PostId, like.UserId, like.Value).
		WillReturnError(sqlErr)
	s.Mock.ExpectRollback()

	res, err := s.repo.Add(like)

	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_CloseRowError() {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(pq.ErrNotSupported)
	expRes := int64(app.InvalidInt)

	like := TestLike(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
		WithArgs(like.PostId, like.UserId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{""}).CloseError(sqlErr))

	s.Mock.ExpectRollback()

	res, err := s.repo.Add(like)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_UpdateError() {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`
	queryUpdate := `UPDATE posts SET likes = likes + $2 WHERE posts_id = $1 RETURNING likes;`

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(pq.ErrNotSupported)
	expRes := int64(app.InvalidInt)

	like := TestLike(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
		WithArgs(like.PostId, like.UserId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, like.Value).
		WillReturnError(sqlErr)

	s.Mock.ExpectRollback()

	res, err := s.repo.Add(like)

	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_CommitError() {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`
	queryUpdate := `UPDATE posts SET likes = likes + $2 WHERE posts_id = $1 RETURNING likes;`

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(pq.ErrNotSupported)
	expRes := int64(app.InvalidInt)
	expLikes := 1

	like := TestLike(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
		WithArgs(like.PostId, like.UserId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).AddRow(expLikes)).
		RowsWillBeClosed()

	s.Mock.ExpectCommit().WillReturnError(sqlErr)

	res, err := s.repo.Add(like)

	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryAdd_OK() {
	queryInsert := `INSERT INTO likes (post_id, users_id, value) VALUES ($1, $2, $3 > 0)`
	queryUpdate := `UPDATE posts SET likes = likes + $2 WHERE posts_id = $1 RETURNING likes;`

	expLikes := int64(1)
	expRes := expLikes

	like := TestLike(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
		WithArgs(like.PostId, like.UserId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{})).
		RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, like.Value).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expLikes)).RowsWillBeClosed()

	s.Mock.ExpectCommit()

	res, err := s.repo.Add(like)

	assert.Equal(s.T(), expRes, res)
	assert.NoError(s.T(), err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryDelete_NotFoundLike() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`

	like := TestLike(s.T())
	expValue := false
	sqlErr := sql.ErrNoRows
	expErr := repository.NotFound
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expValue)).
		WillReturnError(sqlErr)

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryDelete_InternalError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`

	like := TestLike(s.T())
	expValue := false
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expValue)).
		WillReturnError(sqlErr)

	res, err := s.repo.Delete(like.ID)

	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteLikesRepository) TestLikesRepositoryDelete_BeginTransactionError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`

	like := TestLike(s.T())
	expValue := false
	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expValue))

	s.Mock.ExpectBegin().WillReturnError(sqlErr)

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}
func (s *SuiteLikesRepository) TestLikesRepositoryDelete_UpdatePostsError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expRes := int64(app.InvalidInt)
	countLikes := 1

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(countLikes)).WillReturnError(sqlErr)
	s.Mock.ExpectRollback()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}

func (s *SuiteLikesRepository) TestLikesRepositoryDelete_UpdatePostsRowCloseError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expCountLikes := 1
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes).CloseError(sqlErr))
	s.Mock.ExpectRollback()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}

func (s *SuiteLikesRepository) TestLikesRepositoryDelete_DeleteLikeError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expCountLikes := 1
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes)).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow()).
		WillReturnError(sqlErr)

	s.Mock.ExpectRollback()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}

func (s *SuiteLikesRepository) TestLikesRepositoryDelete_DeleteLikeCloseRowError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expCountLikes := 1
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes)).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			AddRow().
			CloseError(sqlErr))
	s.Mock.ExpectRollback()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}

func (s *SuiteLikesRepository) TestLikesRepositoryDelete_CommitError() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1

	sqlErr := pq.ErrNotSupported
	expErr := repository.NewDBError(sqlErr)
	expCountLikes := 1
	expRes := int64(app.InvalidInt)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes)).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			AddRow()).RowsWillBeClosed()
	s.Mock.ExpectCommit().WillReturnError(sqlErr)

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.Equal(s.T(), expErr, err)

}
func (s *SuiteLikesRepository) TestLikesRepositoryDelete_OK_DeleteLike() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	like := TestLike(s.T())
	expBoolValue := true
	convertToLike := 1
	expCountLikes := 1
	expRes := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes)).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			AddRow()).RowsWillBeClosed()
	s.Mock.ExpectCommit()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.NoError(s.T(), err)

}

func (s *SuiteLikesRepository) TestLikesRepositoryDelete_OK_DeleteDislike() {
	querySelect := `SELECT post_id, value FROM likes WHERE likes_id = $1`
	queryUpdate := `UPDATE posts SET likes = likes - $2 WHERE posts_id = $1 RETURNING likes;`
	queryDelete := `DELETE FROM likes WHERE likes_id = $1;`

	like := TestLike(s.T())
	expBoolValue := false
	convertToLike := -1
	expCountLikes := 1
	expRes := int64(1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(querySelect)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "value"}).
			AddRow(like.PostId, expBoolValue))

	s.Mock.ExpectBegin()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryUpdate)).
		WithArgs(like.PostId, convertToLike).
		WillReturnRows(sqlmock.NewRows([]string{"likes"}).
			AddRow(expCountLikes)).RowsWillBeClosed()

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryDelete)).
		WithArgs(like.ID).
		WillReturnRows(sqlmock.NewRows([]string{}).
			AddRow()).RowsWillBeClosed()
	s.Mock.ExpectCommit()

	res, err := s.repo.Delete(like.ID)
	assert.Equal(s.T(), expRes, res)
	assert.NoError(s.T(), err)

}
func TestLikesRepository(t *testing.T) {
	suite.Run(t, new(SuiteLikesRepository))
}
