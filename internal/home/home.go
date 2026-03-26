package home

import (
	"os"
	"path/filepath"

	"github.com/0ghny/gitconfig/internal/config"
	log "github.com/sirupsen/logrus"
)

// Returns the application home directory inside users home
// or current directory
// GetAppHomedir returns the path <userHome>/<appname>. Falls back to "./"
// if the user home cannot be determined.
func GetAppHomedir(appname string) string {
	appHomeDir, err := os.UserHomeDir()
	if err != nil {
		appHomeDir = "./"
	}
	return filepath.Join(appHomeDir, appname)
}

// GetConfigsDir returns the path to the gitconfigs directory without creating it.
// Safe to call from domain code: pure path computation, no filesystem side effects.
func GetConfigsDir() string {
	return GetAppHomedir(config.HomeDirName)
}

// ensure application home directory exists
// EnsureHome creates the application data directory if it does not exist
// and returns its path.
func EnsureHome() (string, error) {
	homeDir := GetAppHomedir(config.HomeDirName)
	err := os.MkdirAll(homeDir, 0755)

	if err != nil {
		return "", err
	} else {
		return homeDir, nil
	}
}

// GetHome ensures the gitconfigs directory exists and returns its path.
// GetHome is equivalent to EnsureHome but terminates the process on error.
// Prefer EnsureHome where error handling is possible.
func GetHome() string {
	dir, err := EnsureHome()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// https://freshman.tech/snippets/go/home-directory/
// GetUserHome returns the current user's home directory.
func GetUserHome() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir, nil
}
