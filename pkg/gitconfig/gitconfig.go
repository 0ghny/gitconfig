package gitconfig

import (
	"path/filepath"

	"github.com/0ghny/gitconfig/internal/home"
	"github.com/0ghny/go-libx/pkg/iox"
)

const (
	gitConfigFileName string = ".gitconfig"
)

func GetUserGitConfigPath() string {
	userHomeDir, err := home.GetUserHome()
	if err != nil {
		panic(err)
	}
	return filepath.Join(userHomeDir, gitConfigFileName)
}

func Exists(path string) bool {
	return iox.FileExists(path)
}
