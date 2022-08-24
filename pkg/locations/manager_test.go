package locations

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/0ghny/gitconfig/internal/filesystem"
	"github.com/0ghny/gitconfig/pkg/gitconfig"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// This variable content should be CLOSEST to the file line start, since
	// our regex is ^ otherwise this test will fail
	gitConfigContent = `
# ..................................................................................
# USER default settings
# ..................................................................................
[user]
	name = unknown
	email = __redacted_to_stop_spam_scrapers__@localhost.local

# ..................................................................................
# APPLY default settings
# ..................................................................................
[apply]
	# Detect whitespace errors when applying a patch
	whitespace = fix
# ..................................................................................
# CREDENTIAL default settings
# ..................................................................................
[credential]
	# leaving this empty, disable the credentials manager
	helper =
# ..................................................................................
# HELP default settings
# ..................................................................................
[help]
	# Automatically correct and execute mistyped commands
	autocorrect = 1
# ..................................................................................
# PUSH default settings
# ..................................................................................
[push]
	default = simple
# ..................................................................................
# COMMIT default settings
# ..................................................................................
[commit]
	gpgsign = false
# ..................................................................................
# Specific folders settings
# ..................................................................................
# gitconfig.location.key location
[includeIf "gitdir:~/tmp/location/"]
	path = ~/.gitconfigs/location.gitconfig

# gitconfig.location.key location1
[includeIf "gitdir:~/location1/"]
	path = ~/.gitconfigs/location1.gitconfig

# gitconfig.location.key location2
[includeIf "gitdir:/var/lib/locations2/"]
	path = ~/.gitconfigs/location2.gitconfig`
)

var (
	aFS *afero.Afero
)

func newMockLocationManager(fileContent string) *LocationManager {
	gitConfigPath := gitconfig.GetUserGitConfigPath()
	aFS = filesystem.NewMemFs()
	aFS.WriteFile(gitConfigPath, []byte(fileContent), fs.ModeAppend)
	return NewLocationManager(gitConfigPath, aFS)
}

func TestGetLocations_WithValidFile_ShouldReturnsLocations(t *testing.T) {
	locationMgr := newMockLocationManager(gitConfigContent)
	locations, err := locationMgr.GetLocations()

	assert.Nil(t, err)
	assert.Equal(t, len(locations), 3)
}

func TestGetLocations_WithNoLocationInFile_ShouldReturnsAnEmptyArrayOfLocations(t *testing.T) {
	locationMgr := newMockLocationManager(`
# gitconfig.location.key location
[includeIf "gitdir:~/tmp/location/"]
	path = ~/.gitconfigs/location.gitconfig

# invalid badspelled header location1
[includeIf "gitdir:~/location1/"]
	path = ~/.gitconfigs/location1.gitconfig

# gitconfig.location.key location2
[includeIf "gitdir:/var/lib/locations2/"]
	path = ~/.gitconfigs/location2.gitconfig
`)
	locations, err := locationMgr.GetLocations()

	assert.Nil(t, err)
	assert.Equal(t, len(locations), 2)
}

func TestGetLocations_WithSomeInvalidLocations_ShouldReturnsOnlyValidLocations(t *testing.T) {
	locationMgr := newMockLocationManager(``)
	locations, err := locationMgr.GetLocations()

	assert.Nil(t, err)
	assert.Equal(t, len(locations), 0)
}

func TestFindLocationByKey(t *testing.T) {
	locationMgr := newMockLocationManager(gitConfigContent)
	locationKey := "location1"
	l, err := locationMgr.FindLocationByKey(locationKey)

	assert.Nil(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, l.Key, locationKey)
	assert.Equal(t, l.Path, "~/location1/")
	assert.Equal(t, l.ConfigFile, "~/.gitconfigs/location1.gitconfig")
}

func TestFindLocationByKey_WithNonExistingKey_ShouldReturnsNil(t *testing.T) {
	locationMgr := newMockLocationManager(gitConfigContent)
	locationKey := "location999"
	l, err := locationMgr.FindLocationByKey(locationKey)

	assert.Nil(t, l)
	assert.Nil(t, err)
}

func TestFindLocationByKey_WithInvalidFile_ShouldReturnsError(t *testing.T) {
	AFS := filesystem.NewMemFs()
	locationMgr := NewLocationManager("/non/existent/path", AFS)
	locationKey := "location1"
	l, err := locationMgr.FindLocationByKey(locationKey)

	assert.Nil(t, l)
	assert.NotNil(t, err)
}

func TestSaveLocation_WithValidLocation_ShouldAddLocationToGitConfigAndCreateLocationConfigFile(t *testing.T) {
	locationMgr := newMockLocationManager(gitConfigContent)
	key := "aKey"
	path := "/tmp/newlocation"
	// Save new location, it should
	//  1. Add section to gitconfig
	//  2. Create file from template templates.go into configured gitconfig home with key as name
	err := locationMgr.SaveLocation(key, path)
	require.Nil(t, err)
	// Get the just created location from file, to check it was created successfully
	l, err := locationMgr.FindLocationByKey(key)
	require.Nil(t, err)
	require.NotNil(t, l)
	// Check if location config file has been created
	_, err = aFS.Stat(l.ConfigFile)
	require.False(t, errors.Is(err, os.ErrNotExist))
	// Check permissions on file, should be 0664
	//currentFileMode := fileInfo.Mode()
	//expectedFileMode := fs.FileMode(int(0664))
	//require.True(t, currentFileMode == expectedFileMode)
}
