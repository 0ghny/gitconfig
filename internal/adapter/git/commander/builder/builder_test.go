package builder_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/0ghny/gitconfig/internal/adapter/git/commander/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommandBuilder_Build_CreatesValidCommand(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").Build()

	require.NoError(t, err)
	assert.NotNil(t, cmd)
}

func TestCommandBuilder_WithArguments_SetsArgs(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").
		WithArguments("hello", "world").
		Build()

	require.NoError(t, err)
	assert.Contains(t, cmd.Args, "hello")
	assert.Contains(t, cmd.Args, "world")
}

func TestCommandBuilder_WithEnvironmentVariable_IncludesVarAndPreservesOSEnv(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").
		WithEnvironmentVariable("MY_KEY", "MY_VALUE").
		Build()

	require.NoError(t, err)
	// Custom var must be present
	assert.Contains(t, cmd.Env, "MY_KEY=MY_VALUE")
	// OS environment must also be preserved (PATH is always set)
	hasPath := false
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "PATH=") {
			hasPath = true
			break
		}
	}
	assert.True(t, hasPath, "OS PATH variable should be inherited")
}

func TestCommandBuilder_WithNoEnvVars_EnvIsNil(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").Build()

	require.NoError(t, err)
	// nil Env tells exec.Cmd to inherit the process environment automatically
	assert.Nil(t, cmd.Env)
}

func TestCommandBuilder_WithTimeout_BuildsCommandSuccessfully(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").
		WithArguments("hi").
		WithTimeout(5 * time.Second).
		Build()

	require.NoError(t, err)
	assert.Contains(t, cmd.Args, "hi")
}

func TestCommandBuilder_WithEnvironmentVariables_AllVarsPresent(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").
		WithEnvironmentVariables([]string{"A=1", "B=2"}).
		Build()

	require.NoError(t, err)
	assert.Contains(t, cmd.Env, "A=1")
	assert.Contains(t, cmd.Env, "B=2")
}

func TestCommandBuilder_IncludeOSEnvironment_MergesWithAdditionalVars(t *testing.T) {
	cmd, err := builder.NewCommandBuilder("echo").
		IncludeOSEnvironment().
		WithEnvironmentVariable("EXTRA", "VALUE").
		Build()

	require.NoError(t, err)
	// The custom var must be present
	assert.Contains(t, cmd.Env, "EXTRA=VALUE")
	// OS env is also present
	assert.Greater(t, len(cmd.Env), len(os.Environ()))
}
