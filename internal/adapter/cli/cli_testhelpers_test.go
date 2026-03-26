package cli_test

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/cli"
	"github.com/0ghny/gitconfig/internal/adapter/filesystem"
	"github.com/0ghny/gitconfig/internal/application/locations"
	"github.com/0ghny/gitconfig/internal/domain/location"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// gitconfig file with two pre-configured locations used across most tests.
const testGitConfig = `[user]
	name = test user
	email = test@localhost
# gitconfig.location.key location2
[includeIf "gitdir:~/location2/"]
	path = ~/.gitconfigs/location2.gitconfig
# gitconfig.location.key location1
[includeIf "gitdir:~/location1/"]
	path = ~/.gitconfigs/location1.gitconfig
`

const testGitConfigPath = "/tmp/test-gitconfig"

// testEnv holds the shared state for a single test scenario.
type testEnv struct {
	gitconfigPath string
	afs           *afero.Afero
	factory       cli.LocationServiceFactory
}

// newTestEnv creates an in-memory filesystem pre-loaded with the given gitconfig content.
func newTestEnv(t *testing.T, gitconfigContent string) *testEnv {
	t.Helper()
	afs := filesystem.NewMemFs()
	if gitconfigContent != "" {
		require.NoError(t, afs.WriteFile(testGitConfigPath, []byte(gitconfigContent), fs.FileMode(0644)))
	}
	factory := func(path string) location.Service {
		return locations.NewLocationManager(path, afs)
	}
	return &testEnv{gitconfigPath: testGitConfigPath, afs: afs, factory: factory}
}

// run executes a fresh cobra command tree with the given args and optional stdin content.
// The --git-config flag is always injected so commands target the test file.
func (e *testEnv) run(t *testing.T, stdin string, args ...string) (output string, err error) {
	t.Helper()
	out := &bytes.Buffer{}
	root := cli.RootCmdWithDeps("test", e.factory)
	root.SetOut(out)
	root.SetErr(out)
	if stdin != "" {
		root.SetIn(strings.NewReader(stdin))
	}
	allArgs := append([]string{"--git-config", e.gitconfigPath}, args...)
	root.SetArgs(allArgs)
	root.TraverseChildren = true
	_, err = root.ExecuteC()
	return out.String(), err
}

// preCreateConfigFile creates a placeholder location config file in the in-memory FS.
// Required for delete tests because DeleteLocation calls fs.Remove on the config file.
func (e *testEnv) preCreateConfigFile(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, e.afs.WriteFile(path, []byte(""), fs.FileMode(0644)))
}
