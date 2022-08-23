package commander

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

func RunOutput(command string, args ...string) (output string, err error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, command, args...)
	return RunCommandCombined(*cmd)
}

func RunInteractive(command string, args ...string) (err error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, command, args...)
	return RunCommandInteractive(*cmd)
}

func RunCommandCombined(cmd exec.Cmd) (output string, err error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}

	return string(out), nil
}

func RunCommandInteractive(cmd exec.Cmd) (err error) {
	// Connectes all in/out to the process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// Run command
	return cmd.Run()
}
