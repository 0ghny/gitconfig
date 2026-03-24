package gitconfig

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/0ghny/gitconfig/internal/adapter/filesystem"
	"github.com/0ghny/gitconfig/internal/home"
	"github.com/spf13/afero"
)

const (
	gitConfigFileName string      = ".gitconfig"
	filePerms         fs.FileMode = fs.FileMode(int(0644))
)

// GitConfig manages read and write operations on a single .gitconfig file.
type GitConfig struct {
	path string
	fs   *afero.Afero
}

// NewGitConfig creates a GitConfig for the given path and filesystem.
// If path is empty, the user's default ~/.gitconfig is used.
// If fs is nil, the OS filesystem is used.
func NewGitConfig(path string, fs *afero.Afero) *GitConfig {
	if path == "" {
		path = GetUserGitConfigPath()
	}
	if fs == nil {
		fs = filesystem.NewOsFs()
	}
	return &GitConfig{path: path, fs: fs}
}

// GetUserGitConfigPath returns the path to the user's default ~/.gitconfig file.
func GetUserGitConfigPath() string {
	userHomeDir, err := home.GetUserHome()
	if err != nil {
		panic(err)
	}
	return filepath.Join(userHomeDir, gitConfigFileName)
}

// Exists reports whether the file at path exists on the OS filesystem.
func Exists(path string) bool {
	return filesystem.FileExists(path)
}

// GetContent returns the full text content of the gitconfig file.
func (g *GitConfig) GetContent() (string, error) {
	fileBytes, err := g.fs.ReadFile(g.path)
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

// AppendSection appends a new section to the end of the gitconfig file.
func (g *GitConfig) AppendSection(section string) error {
	f, err := g.fs.OpenFile(g.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, filePerms)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(section)
	return err
}

// WriteContent overwrites the entire gitconfig file with the given content.
func (g *GitConfig) WriteContent(content string) error {
	return g.fs.WriteFile(g.path, []byte(content), filePerms)
}
