package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configService, err := NewConfigService("test", "./config_test.yaml")
	assert.Nil(t, err)

	s1, ok := configService.Get("test.config.string")
	assert.True(t, ok)
	assert.Equal(t, "string", s1.(string))

	s2, ok := configService.GetString("test.config.string")
	assert.True(t, ok)
	assert.Equal(t, "string", s2)
}
