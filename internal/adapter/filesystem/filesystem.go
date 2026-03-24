package filesystem

import (
	"github.com/spf13/afero"
)

// NewOsFs returns an afero.Afero backed by the real OS filesystem.
func NewOsFs() *afero.Afero {
	FS := afero.NewOsFs()
	return &afero.Afero{Fs: FS}
}

// NewMemFs returns an afero.Afero backed by an in-memory filesystem.
// Primarily used in tests to avoid touching the real filesystem.
func NewMemFs() *afero.Afero {
	FS := afero.NewMemMapFs()
	return &afero.Afero{Fs: FS}
}
