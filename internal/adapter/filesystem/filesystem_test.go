package filesystem_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileExists_WithExistingFile_ReturnsTrue(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "testfile")
	require.NoError(t, err)
	f.Close()

	assert.True(t, filesystem.FileExists(f.Name()))
}

func TestFileExists_WithNonExistentPath_ReturnsFalse(t *testing.T) {
	assert.False(t, filesystem.FileExists("/nonexistent/path/file.txt"))
}

func TestFileExists_WithDirectory_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()

	assert.False(t, filesystem.FileExists(dir))
}

func TestGetwd_ReturnsNonEmptyString(t *testing.T) {
	wd := filesystem.Getwd()

	assert.NotEmpty(t, wd)
}

func TestGetBaseWd_ReturnsNonEmptyBaseName(t *testing.T) {
	base := filesystem.GetBaseWd()

	assert.NotEmpty(t, base)
	// Base should not contain a path separator
	assert.False(t, strings.Contains(base, string(filepath.Separator)))
}

func TestGetAbsWd_ReturnsAbsolutePath(t *testing.T) {
	absPath, err := filesystem.GetAbsWd()

	require.NoError(t, err)
	assert.True(t, filepath.IsAbs(absPath))
}

func TestNewOsFs_ReturnsNonNil(t *testing.T) {
	fs := filesystem.NewOsFs()

	assert.NotNil(t, fs)
}

func TestNewMemFs_ReturnsNonNil(t *testing.T) {
	fs := filesystem.NewMemFs()

	assert.NotNil(t, fs)
}
