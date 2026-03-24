package filesystem

import (
	"os"
)

// Check wether a file exists or not.
// Note: a directory is not a file
func FileExists(f string) bool {
	fileInfo, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}
