package filesystem

import (
	"os"
	"path/filepath"
)

// GetBaseWd returns the base name of the current working directory.
func GetBaseWd() string {
	currDir, err := os.Getwd()
	if err != nil {
		currDir = "./"
	}
	return filepath.Base(currDir)
}

// Getwd returns the current working directory path, falling back to "./" on error.
func Getwd() string {
	currDir, err := os.Getwd()
	if err != nil {
		return "./"
	}
	return currDir
}

// GetAbsWd returns the absolute path of the current working directory.
func GetAbsWd() (string, error) {
	absPath, err := filepath.Abs(Getwd())
	if err != nil {
		return "", err
	}
	return absPath, nil
}
