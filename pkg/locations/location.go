package locations

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/0ghny/gitconfigs/internal/home"
)

const (
	locationSectionTemplate string = `
# gitconfigs.location.key {{.Key}}
[includeIf "gitdir:{{.Path}}/"]
	path = {{.ConfigFile}}
`
)

type Location struct {
	Key  string
	Path string
	// filePath.join(home.GetHome, fmt.Sprintf("%s.gitconfig", Key))
	ConfigFile string
}

func NewLocation(key string, location string) *Location {
	return &Location{
		Key:        key,
		Path:       location,
		ConfigFile: filepath.Join(home.GetHome(), fmt.Sprintf("%s.gitconfig", key)),
	}
}

// Returns the section composed from the location object
// this section is the text that usually goes into the gitconfig file
func (l Location) ToSection() (string, error) {
	var section bytes.Buffer
	sectionTemplate := template.New("SectionTemplate")

	sectionTemplate, err := sectionTemplate.Parse(locationSectionTemplate)
	if err != nil {
		// error creating template
		return "", err
	}

	err = sectionTemplate.Execute(&section, l)
	if err != nil {
		// error parsing template with data
		return "", err
	}

	// returns section
	return section.String(), nil
}
