package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCmd_NoArgs_ReturnsValidationError(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	_, err := env.run(t, "", "config")

	assert.Error(t, err)
}

func TestConfigCmd_TooManyArgs_ReturnsValidationError(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	_, err := env.run(t, "", "config", "key", "value", "extra")

	assert.Error(t, err)
}

func TestConfigCmd_WithKeyAndNonExistingLocation_ReturnsError(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	_, err := env.run(t, "", "config", "--key", "nonexistent", "user.name")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}
