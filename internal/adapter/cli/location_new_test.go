package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationNewCmd_WithValidArgs_AddsLocationAndPrintsSuccess(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	out, err := env.run(t, "", "location", "new", "newloc", "--location", "/code/newproject")

	require.NoError(t, err)
	assert.Contains(t, out, "newloc")
	assert.Contains(t, out, "saved successfully")

	// Verify it persisted by reading it back.
	svc := env.factory(testGitConfigPath)
	l, readErr := svc.FindLocationByKey("newloc")
	require.NoError(t, readErr)
	require.NotNil(t, l)
	assert.Equal(t, "/code/newproject", l.Path)
}

func TestLocationNewCmd_WithDuplicateKey_UpdatesExistingLocation(t *testing.T) {
	env := newTestEnv(t, testGitConfig)
	// Pre-create the work config file so the update rename does not fail.
	env.preCreateConfigFile(t, "~/.gitconfigs/work.gitconfig")

	_, err := env.run(t, "", "location", "new", "work", "--location", "/code/work-updated")

	require.NoError(t, err)
	svc := env.factory(testGitConfigPath)
	l, readErr := svc.FindLocationByKey("work")
	require.NoError(t, readErr)
	require.NotNil(t, l)
	assert.Equal(t, "/code/work-updated", l.Path)
}
