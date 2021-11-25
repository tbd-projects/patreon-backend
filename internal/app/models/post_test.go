package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePost_ValidateIncorrectTitle(t *testing.T) {
	post := TestUpdatePost()
	post.Title = ""
	expErr := EmptyTitle

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestUpdatePost_ValidateIncorrectAwards(t *testing.T) {
	post := TestUpdatePost()
	post.Awards = -2
	expErr := InvalidAwardsId

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestUpdatePost_ValidateOk(t *testing.T) {
	post := TestUpdatePost()

	res := post.Validate()
	assert.NoError(t, res)
}
func TestCreatePost_ValidateInvalidTitle(t *testing.T) {
	post := TestCreatePost()
	post.Title = ""
	expErr := EmptyTitle

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestCreatePost_ValidateInvalidCreatorId(t *testing.T) {
	post := TestCreatePost()
	post.CreatorId = -1
	expErr := InvalidCreatorId

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestCreatePost_ValidateInvalidAwardId(t *testing.T) {
	post := TestCreatePost()
	post.Awards = -2
	expErr := InvalidAwardsId

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestCreatePost_ValidateInvalidACombinatinon(t *testing.T) {
	post := TestCreatePost()
	post.Awards = -2
	post.Title = ""

	res := post.Validate()
	assert.True(t, res == InvalidAwardsId || res == EmptyTitle)
}
func TestCreatePost_Validate_OK(t *testing.T) {
	post := TestCreatePost()

	res := post.Validate()
	assert.NoError(t, res)
}

func TestAttachWithoutLevel_ValidateIncorrectPostId(t *testing.T) {
	post := TestAttachWithoutLevel()
	post.PostId = -1
	expErr := InvalidPostId

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestAttachWithoutLevel_ValidateIncorrectAttachWithoutLevelType(t *testing.T) {
	post := TestAttachWithoutLevel()
	post.Type = "invalid"
	expErr := InvalidType

	res := post.Validate()
	assert.Equal(t, expErr, res)
}
func TestAttachWithoutLevel_ValidateIncorrectAttachWithoutLevelCombination(t *testing.T) {
	post := TestAttachWithoutLevel()
	post.Type = "invalid"
	post.PostId = -2

	res := post.Validate()
	assert.True(t, res == InvalidPostId || res == InvalidType)
}
func TestAttachWithoutLevel_ValidateIncorrectAttachWithoutLevelOk(t *testing.T) {
	post := TestAttachWithoutLevel()
	res := post.Validate()
	assert.NoError(t, res)
}
