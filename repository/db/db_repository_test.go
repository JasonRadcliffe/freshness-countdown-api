package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb_GetDishes(t *testing.T) {
	assert.Equal(t, "", "")
}
func TestDb_GetDishes_QueryError(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestDb_GetDishes_RowScanError(t *testing.T) {
	assert.Equal(t, "", "")
}
