package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIRepositories(t *testing.T) {
	repo := NewAPIRepositories()
	assert.NotNil(t, repo.PackSizes)
}
