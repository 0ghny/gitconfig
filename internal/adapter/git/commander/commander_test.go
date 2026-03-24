package commander_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/0ghny/gitconfig/internal/adapter/git/commander"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunOutput_WithEcho_ReturnsOutput(t *testing.T) {
	out, err := commander.RunOutput("echo", "hello")

	require.NoError(t, err)
	assert.Equal(t, "hello\n", out)
}

func TestRunOutput_WithNonExistentCommand_ReturnsError(t *testing.T) {
	_, err := commander.RunOutput("nonexistentcommand_xyz123")

	assert.Error(t, err)
}

func TestRunCommandCombined_WithSuccessfulCommand_ReturnsOutput(t *testing.T) {
	cmd := *exec.Command("echo", "combined")

	out, err := commander.RunCommandCombined(cmd)

	require.NoError(t, err)
	assert.True(t, strings.Contains(out, "combined"))
}

func TestRunCommandCombined_WithFailingCommand_ReturnsError(t *testing.T) {
	// `ls /path/that/does/not/exist/xyz` exits non-zero
	cmd := *exec.Command("ls", "/path/that/does/not/exist/xyz_abc_123")

	_, err := commander.RunCommandCombined(cmd)

	assert.Error(t, err)
}
