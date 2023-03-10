package stamp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	_, err := NewDefaultConfig()
	assert.NoError(t, err)
}

func TestNewConfig_Valid(t *testing.T) {
	path := filepath.Join("testdata", "config", "valid.yml")
	config, err := NewConfig(path)
	assert.NoError(t, err)
	assert.Equal(t, "~/some/dir", config.StorePath)
	// viper converts all keys to lowercase :shrug:
	assert.Nil(t, config.Defaults["UserName"])
	assert.Equal(t, "test-user", config.Defaults["username"])
}

func TestNewConfig_Invalid(t *testing.T) {
	path := filepath.Join("testdata", "config", "invalid.yml")
	_, err := NewConfig(path)
	assert.ErrorContains(t, err, "cannot parse 'debug' as bool")
}
