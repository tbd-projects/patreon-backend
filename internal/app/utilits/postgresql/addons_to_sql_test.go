package postgresql_utilits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CustomBind(t *testing.T) {
	queryInput := "SELECT n_live_tup FROM (?, ?, ?)"
	queryOutput := "SELECT n_live_tup FROM ($1, $2, $3)"
	startIndex := 1

	res := CustomRebind(startIndex, queryInput)
	assert.Equal(t, res, queryOutput)
}
