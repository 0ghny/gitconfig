package cli_test

import (
	"bytes"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/cli"
	"github.com/0ghny/gitconfig/internal/domain/location"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd_VersionFlag_PrintsVersion(t *testing.T) {
	env := newTestEnv(t, "")

	out, err := env.run(t, "", "--version")

	require.NoError(t, err)
	assert.Contains(t, out, "test") // version string passed to RootCmdWithDeps
}

func TestRootCmd_UnknownSubcommand_ReturnsError(t *testing.T) {
	// Cannot use env.run() here: its TraverseChildren=true makes cobra silently
	// fall back to the root command instead of returning an error for unknown
	// subcommands. Build the root command directly to match production behaviour.
	out := &bytes.Buffer{}
	factory := func(_ string) location.Service { return nil }
	root := cli.RootCmdWithDeps("test", factory)
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"nonexistent"})

	_, err := root.ExecuteC()

	assert.Error(t, err)
}
