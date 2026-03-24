package commander

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

// RunOutput runs command with the given arguments and returns the combined
// stdout+stderr output as a string.
func RunOutput(command string, args ...string) (output string, err error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, command, args...)
	return RunCommandCombined(*cmd)
}

// RunInteractive runs command with its stdin/stdout/stderr wired to the
// current process, allowing interactive use.
func RunInteractive(command string, args ...string) (err error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, command, args...)
	return RunCommandInteractive(*cmd)
}

// RunCommandCombined executes a pre-built exec.Cmd and returns the combined
// stdout+stderr output. Errors from the process are wrapped with the output text.
func RunCommandCombined(cmd exec.Cmd) (output string, err error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}

	return string(out), nil
}

// RunCommandInteractive wires stdin/stdout/stderr of cmd to the current
// process and runs it, enabling interactive commands.
func RunCommandInteractive(cmd exec.Cmd) (err error) {
	// Connectes all in/out to the process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// Run command
	return cmd.Run()
}
