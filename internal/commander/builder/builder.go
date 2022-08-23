package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

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
	// Creates command based on builder configuration
	var theCmd *exec.Cmd
	if b.timeout > 0 {
		theCmd = b.createCommandWithTimeout()
	} else {
		theCmd = b.createCommand()
	}

	theCmd.Env = b.env
	return *theCmd, nil
}
func (b *CommandBuilder) WithArguments(a ...string) *CommandBuilder {
	b.args = a
	return b
}
func (b *CommandBuilder) WithEnvironmentVariable(key string, value string) *CommandBuilder {
	b.env = append(b.env, fmt.Sprintf("%s=%s", key, value))
	return b
}
func (b *CommandBuilder) WithEnvironmentVariables(envs []string) *CommandBuilder {
	b.env = append(b.env, envs...)
	return b
}
func (b *CommandBuilder) IncludeOSEnvironment() *CommandBuilder {
	b.env = append(b.env, os.Environ()...)
	return b
}
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
