package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationCmd_WithExistingKey_PrintsLocationDetails(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	out, err := env.run(t, "", "location", "work")

	require.NoError(t, err)
	assert.Contains(t, out, "work")
	assert.Contains(t, out, "~/work")
}

func TestLocationCmd_WithNonExistingKey_ReturnsError(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	_, err := env.run(t, "", "location", "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLocationCmd_WithoutKey_ShowsHelp(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	out, err := env.run(t, "", "location")

	require.NoError(t, err)
	assert.Contains(t, out, "Usage")
}
