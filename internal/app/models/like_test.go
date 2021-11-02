package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLike_ValidateInvalidValue(t *testing.T) {
	like := TestLike()
	like.Value = -2
	expErr := InvalidLikeValue
	assert.Equal(t, expErr, like.Validate())
}

func TestLike_ValidateCorrect(t *testing.T) {
	like := TestLike()
	assert.NoError(t, like.Validate())
}
