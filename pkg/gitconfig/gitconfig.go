package gitconfig

import (
	"path/filepath"

	"github.com/0ghny/gitconfigs/internal/home"
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
