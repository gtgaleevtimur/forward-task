package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	r := NewRepository()
	assert.IsType(t, &Repository{}, r)
}

func TestRepository_UpdateRegistry(t *testing.T) {
	r := NewRepository()
	err := r.UpdateRegistry()
	assert.NoError(t, err)
}

func TestRepository_Get(t *testing.T) {
	r := NewRepository()
	value, err := r.Get("0")
	assert.NoError(t, err)
	assert.Equal(t, "", value)
}
