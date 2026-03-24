package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationDeleteCmd_ConfirmedWithY_DeletesLocationAndPrintsSuccess(t *testing.T) {
	env := newTestEnv(t, testGitConfig)
	env.preCreateConfigFile(t, "~/.gitconfigs/work.gitconfig")

	out, err := env.run(t, "y\n", "location", "delete", "work")

	require.NoError(t, err)
	assert.Contains(t, out, "deleted successfully")

	svc := env.factory(testGitConfigPath)
	l, readErr := svc.FindLocationByKey("work")
	require.NoError(t, readErr)
	assert.Nil(t, l, "location should have been removed")
}

func TestLocationDeleteCmd_ConfirmedWithEnter_DeletesLocation(t *testing.T) {
	env := newTestEnv(t, testGitConfig)
	env.preCreateConfigFile(t, "~/.gitconfigs/work.gitconfig")

	out, err := env.run(t, "\n", "location", "delete", "work")

	require.NoError(t, err)
	assert.Contains(t, out, "deleted successfully")
}

func TestLocationDeleteCmd_ConfirmedWithYes_DeletesLocation(t *testing.T) {
	env := newTestEnv(t, testGitConfig)
	env.preCreateConfigFile(t, "~/.gitconfigs/work.gitconfig")

	out, err := env.run(t, "yes\n", "location", "delete", "work")

	require.NoError(t, err)
	assert.Contains(t, out, "deleted successfully")
}

func TestLocationDeleteCmd_CancelledWithN_DoesNotDeleteAndPrintsCancelled(t *testing.T) {
	env := newTestEnv(t, testGitConfig)
	env.preCreateConfigFile(t, "~/.gitconfigs/work.gitconfig")

	out, err := env.run(t, "n\n", "location", "delete", "work")

	require.NoError(t, err)
	assert.Contains(t, out, "cancelled")

	// Location must still exist.
	svc := env.factory(testGitConfigPath)
	l, readErr := svc.FindLocationByKey("work")
	require.NoError(t, readErr)
	assert.NotNil(t, l)
}

func TestLocationDeleteCmd_WithNonExistingKey_ReturnsErrorBeforePrompt(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	out, err := env.run(t, "", "location", "delete", "ghost")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	// Confirmation prompt must NOT have been shown.
	assert.NotContains(t, out, "Are you sure")
}
