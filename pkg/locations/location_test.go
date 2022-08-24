package locations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSection_ShouldReturnsValidSectionFromLocationObject(t *testing.T) {
	expectedSection := `
# gitconfig.location.key aKey
[includeIf "gitdir:/a/path"]
	path = /configs/aKey.gitconfig
`

	location := Location{
		Key:        "aKey",
		Path:       "/a/path",
		ConfigFile: "/configs/aKey.gitconfig",
	}

	sectionText, err := location.ToSection()

	assert.Nil(t, err)
	assert.NotNil(t, location)
	assert.NotEmpty(t, sectionText)
	assert.Equal(t, sectionText, expectedSection)
}
