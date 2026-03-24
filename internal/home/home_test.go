package home_test

import (
	"os"
	"strings"
	"testing"

	"github.com/0ghny/gitconfig/internal/home"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAppHomedir_ContainsGivenAppName(t *testing.T) {
	dir := home.GetAppHomedir("myapp")

	assert.True(t, strings.HasSuffix(dir, "/myapp"), "expected dir to end with /myapp, got: %s", dir)
}

func TestGetConfigsDir_ContainsDotGitconfigs(t *testing.T) {
	dir := home.GetConfigsDir()

	assert.True(t, strings.HasSuffix(dir, "/.gitconfigs"), "expected dir to end with /.gitconfigs, got: %s", dir)
}

func TestGetUserHome_ReturnsNonEmptyPath(t *testing.T) {
	h, err := home.GetUserHome()

	require.NoError(t, err)
	assert.NotEmpty(t, h)
}

func TestEnsureHome_CreatesDirectoryAndReturnsPath(t *testing.T) {
	dir, err := home.EnsureHome()

	require.NoError(t, err)
	assert.NotEmpty(t, dir)

	info, statErr := os.Stat(dir)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())
}
