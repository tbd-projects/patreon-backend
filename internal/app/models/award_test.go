package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAward_ValidateIncorrectName(t *testing.T) {
	aw := TestAward()
	aw.Name = ""
	expErr := EmptyName
	assert.Equal(t, expErr, aw.Validate())
}
func TestAward_ValidateIncorrectPrice(t *testing.T) {
	aw := TestAward()
	aw.Price = -1
	expErr := IncorrectAwardsPrice
	assert.Equal(t, expErr, aw.Validate())
}
func TestAward_Validate_OK(t *testing.T) {
	aw := TestAward()
	assert.NoError(t, aw.Validate())
}
func TestAward_ValidateIncorrectPriceAndName(t *testing.T) {
	aw := TestAward()
	aw.Price = -1
	aw.Name = ""
	err := aw.Validate()
	assert.True(t, err == IncorrectAwardsPrice || err == EmptyName)
}
