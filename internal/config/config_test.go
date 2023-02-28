package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	conf := NewConfig()
	assert.IsType(t, &Config{}, conf)
}
