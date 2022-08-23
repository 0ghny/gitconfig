package home

import (
	"os"

	"github.com/0ghny/gitconfigs/internal/config"
	"github.com/0ghny/go-libx/pkg/iox"
	log "github.com/sirupsen/logrus"
)

// ensure application home directory exists
func EnsureHome() (string, error) {
	homeDir := iox.GetAppHomedir(config.HomeDirName)
	err := os.MkdirAll(homeDir, 0755)

	if err != nil {
		return "", err
	} else {
		return homeDir, nil
	}
}

// returns the path to the application home directory
func GetHome() string {
	dir, err := EnsureHome()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return dir
}

// https://freshman.tech/snippets/go/home-directory/
func GetUserHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir, nil
}
