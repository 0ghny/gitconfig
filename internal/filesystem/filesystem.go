package filesystem

import (
	"github.com/spf13/afero"
)

func NewOsFs() *afero.Afero {
	FS := afero.NewOsFs()
	return &afero.Afero{Fs: FS}
}

func NewMemFs() *afero.Afero {
	FS := afero.NewMemMapFs()
	return &afero.Afero{Fs: FS}
}
