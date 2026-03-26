package gitconfig

import (
	"io/fs"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPath    = "/tmp/.gitconfig"
	testContent = "[user]\n\tname = test\n\temail = test@localhost\n"
)

func newMockGitConfig(content string) *GitConfig {
	afs := filesystem.NewMemFs()
	_ = afs.WriteFile(testPath, []byte(content), fs.ModeAppend)
	return NewGitConfig(testPath, afs)
}

// NewGitConfig

func TestNewGitConfig_WithEmptyPath_ShouldUseUserDefaultmPath(t *testing.T) {
	gc := NewGitConfig("", nil)
	assert.Equal(t, GetUserGitConfigPath(), gc.path)
}

func TestNewGitConfig_WithCustomPath_ShouldUseProvidedPath(t *testing.T) {
	gc := NewGitConfig(testPath, filesystem.NewMemFs())
	assert.Equal(t, testPath, gc.path)
}

// GetContent

func TestGetContent_WithExistingFile_ShouldReturnContent(t *testing.T) {
	gc := newMockGitConfig(testContent)

	content, err := gc.GetContent()

	require.Nil(t, err)
	assert.Equal(t, testContent, content)
}

func TestGetContent_WithEmptyFile_ShouldReturnEmptyString(t *testing.T) {
	gc := newMockGitConfig("")

	content, err := gc.GetContent()

	require.Nil(t, err)
	assert.Equal(t, "", content)
}

func TestGetContent_WithNonExistentFile_ShouldReturnError(t *testing.T) {
	gc := NewGitConfig("/non/existent/.gitconfig", filesystem.NewMemFs())

	content, err := gc.GetContent()

	assert.NotNil(t, err)
	assert.Equal(t, "", content)
}

// AppendSection

func TestAppendSection_OnExistingFile_ShouldAppendContent(t *testing.T) {
	gc := newMockGitConfig(testContent)
	section := "\n[core]\n\tautocrlf = false\n"

	err := gc.AppendSection(section)

	require.Nil(t, err)
	content, err := gc.GetContent()
	require.Nil(t, err)
	assert.Equal(t, testContent+section, content)
}

func TestAppendSection_OnNonExistentFile_ShouldCreateFileAndWrite(t *testing.T) {
	afs := filesystem.NewMemFs()
	gc := NewGitConfig(testPath, afs)
	section := "[user]\n\tname = new\n"

	err := gc.AppendSection(section)

	require.Nil(t, err)
	content, err := gc.GetContent()
	require.Nil(t, err)
	assert.Equal(t, section, content)
}

func TestAppendSection_CalledTwice_ShouldAccumulateContent(t *testing.T) {
	gc := newMockGitConfig("")
	first := "[user]\n\tname = a\n"
	second := "[core]\n\tbare = false\n"

	require.Nil(t, gc.AppendSection(first))
	require.Nil(t, gc.AppendSection(second))

	content, err := gc.GetContent()
	require.Nil(t, err)
	assert.Equal(t, first+second, content)
}

// WriteContent

func TestWriteContent_ShouldOverwriteExistingContent(t *testing.T) {
	gc := newMockGitConfig(testContent)
	newContent := "[core]\n\tbare = true\n"

	err := gc.WriteContent(newContent)

	require.Nil(t, err)
	content, err := gc.GetContent()
	require.Nil(t, err)
	assert.Equal(t, newContent, content)
}

func TestWriteContent_WithEmptyString_ShouldTruncateFile(t *testing.T) {
	gc := newMockGitConfig(testContent)

	err := gc.WriteContent("")

	require.Nil(t, err)
	content, err := gc.GetContent()
	require.Nil(t, err)
	assert.Equal(t, "", content)
}
