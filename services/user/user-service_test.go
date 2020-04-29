package user

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func TestUser_GetByID(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestUser_GetByEmail(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestUser_Create(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestUser_GenerateTempMatch(t *testing.T) {
	assert.Equal(t, "", "")
}

func TestGetByEmail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.Equal(t, "", "")

	test := 2
	fmt.Println("test", test)

}
