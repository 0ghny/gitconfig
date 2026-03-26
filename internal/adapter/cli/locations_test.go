package cli_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocationsCmd_WithTwoLocations_PrintsBothKeys(t *testing.T) {
	env := newTestEnv(t, testGitConfig)

	out, err := env.run(t, "", "locations")

	require.NoError(t, err)
	assert.Contains(t, out, "location2")
	assert.Contains(t, out, "location1")
}

func TestLocationsCmd_WithNoLocations_PrintsEmptyTableWithoutError(t *testing.T) {
	env := newTestEnv(t, "[user]\n\tname = test\n")

	out, err := env.run(t, "", "locations")

	require.NoError(t, err)
	// Table header must still be present; no location rows.
	assert.Contains(t, out, "KEY")
	assert.NotContains(t, out, "location2")
}

func TestLocationsCmd_WithMissingGitConfigFile_ReturnsError(t *testing.T) {
	env := newTestEnv(t, "") // no file written

	_, err := env.run(t, "", "locations")

	assert.Error(t, err)
}

func TestLocationsCmd_GitConfigFlagIsUsed_ReadsFromSpecifiedFile(t *testing.T) {
	env := newTestEnv(t, "") // no content at default path

	customPath := "/custom/.gitconfig"
	customContent := `# gitconfig.location.key custom
[includeIf "gitdir:~/custom/"]
	path = ~/.gitconfigs/custom.gitconfig
`
	require.NoError(t, env.afs.WriteFile(customPath, []byte(customContent), fs.FileMode(0644)))
	env.gitconfigPath = customPath

	out, err := env.run(t, "", "locations")

	require.NoError(t, err)
	assert.Contains(t, out, "custom")
}
