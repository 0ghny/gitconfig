package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// CommandBuilder accumulates the configuration for a single external command.
type CommandBuilder struct {
	command string
	args    []string
	env     []string
	timeout time.Duration
}

// Get a new command builder instance
func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command: command,
	}
}

// Builds a Command
func (b *CommandBuilder) Build() (exec.Cmd, error) {
	var theCmd *exec.Cmd
	if b.timeout > 0 {
		theCmd = b.createCommandWithTimeout()
	} else {
		theCmd = b.createCommand()
	}

	if len(b.env) > 0 {
		// Always inherit the OS environment so that PATH, HOME and other
		// critical variables remain available. Custom vars are appended last
		// so they take precedence on name collision.
		theCmd.Env = append(os.Environ(), b.env...)
	}
	// When b.env is empty, theCmd.Env stays nil which inherits the process
	// environment automatically — no explicit copy needed.
	return *theCmd, nil
}

// WithArguments sets the positional arguments for the command.
func (b *CommandBuilder) WithArguments(a ...string) *CommandBuilder {
	b.args = a
	return b
}

// WithEnvironmentVariable appends a single KEY=VALUE pair to the command's environment.
func (b *CommandBuilder) WithEnvironmentVariable(key string, value string) *CommandBuilder {
	b.env = append(b.env, fmt.Sprintf("%s=%s", key, value))
	return b
}

// WithEnvironmentVariables appends multiple KEY=VALUE pairs to the command's environment.
func (b *CommandBuilder) WithEnvironmentVariables(envs []string) *CommandBuilder {
	b.env = append(b.env, envs...)
	return b
}

// IncludeOSEnvironment copies the current process's environment into the
// command's environment, so it merges with any additional variables added via
// WithEnvironmentVariable.
func (b *CommandBuilder) IncludeOSEnvironment() *CommandBuilder {
	b.env = append(b.env, os.Environ()...)
	return b
}

// WithTimeout sets a deadline on the command's context. The command is
// cancelled if it has not finished within the given duration.
func (b *CommandBuilder) WithTimeout(t time.Duration) *CommandBuilder {
	b.timeout = t
	return b
}

func (b *CommandBuilder) createCommand() *exec.Cmd {
	ctx := context.Background()
	return exec.CommandContext(ctx, b.command, b.args...)
}
func (b *CommandBuilder) createCommandWithTimeout() *exec.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	return exec.CommandContext(ctx, b.command, b.args...)
}
