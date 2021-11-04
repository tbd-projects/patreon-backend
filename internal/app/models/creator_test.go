package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreator_ValidateIncorrectNickname(t *testing.T) {
	cr := TestCreator()
	cr.Nickname = ""
	expErr := IncorrectCreatorNickname

	err := cr.Validate()

	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}
func TestCreator_ValidateIncorrectCategory(t *testing.T) {
	cr := TestCreator()
	cr.Category = ""
	expErr := IncorrectCreatorCategory

	err := cr.Validate()

	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}
func TestCreator_ValidateIncorrectDescription(t *testing.T) {
	cr := TestCreator()
	cr.Description = ""
	expErr := IncorrectCreatorDescription

	err := cr.Validate()

	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}
func TestCreator_ValidateCombineIncorrectField(t *testing.T) {
	cr := TestCreator()
	cr.Nickname = ""
	cr.Description = ""

	err := cr.Validate()

	assert.Error(t, err)
	assert.True(t, err == IncorrectCreatorDescription || err == IncorrectCreatorNickname)
}
func TestCreator_Validate_OK(t *testing.T) {
	cr := TestCreator()
	err := cr.Validate()

	assert.NoError(t, err)
}
func TestCreator_String(t *testing.T) {
	cr := TestCreator()
	assert.Equal(t, fmt.Sprintf("%v", cr), cr.String())
}
