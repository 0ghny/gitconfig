package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIfGitNotFound(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not found in PATH")
	}
}

func TestGitConfigGet_WithNonExistentFile_ReturnsError(t *testing.T) {
	_, err := git.GitConfigGet("user.name", "/nonexistent/path/.gitconfig")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "/nonexistent/path/.gitconfig")
}

func TestGitConfigSet_WithNonExistentFile_ReturnsError(t *testing.T) {
	err := git.GitConfigSet("user.name", "test", "/nonexistent/path/.gitconfig")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "/nonexistent/path/.gitconfig")
}

func TestGitConfigGet_WithExistingKey_ReturnsValue(t *testing.T) {
	skipIfGitNotFound(t)

	cfgPath := filepath.Join(t.TempDir(), ".gitconfig")
	require.NoError(t, os.WriteFile(cfgPath, []byte("[user]\n\tname = testuser\n"), 0644))

	value, err := git.GitConfigGet("user.name", cfgPath)

	require.NoError(t, err)
	assert.Equal(t, "testuser\n", value)
}

func TestGitConfigSet_WritesKeyToFile(t *testing.T) {
	skipIfGitNotFound(t)

	cfgPath := filepath.Join(t.TempDir(), ".gitconfig")
	require.NoError(t, os.WriteFile(cfgPath, []byte("[user]\n\tname = original\n"), 0644))

	err := git.GitConfigSet("user.name", "updated", cfgPath)

	require.NoError(t, err)
	// Verify the value was written by reading it back
	value, err := git.GitConfigGet("user.name", cfgPath)
	require.NoError(t, err)
	assert.Equal(t, "updated\n", value)
}
