package git

import (
	"fmt"

	"github.com/0ghny/gitconfigs/internal/commander"
	"github.com/0ghny/gitconfigs/internal/commander/builder"
	gcf "github.com/0ghny/gitconfigs/pkg/gitconfig"
)

func GitConfigGet(key string, gitconfig string) (value string, err error) {
	commandBuilder := builder.NewCommandBuilder("git")
	if gitconfig != "" {
		if !gcf.Exists(gitconfig) {
			return "", fmt.Errorf("gitconfig specified not found at %s", gitconfig)
		}
		commandBuilder.WithEnvironmentVariable("GIT_CONFIG", gitconfig)
	}
	gitCommand, err := commandBuilder.
		WithArguments([]string{"config", key}...).
		Build()

	if err != nil {
		// error creating command
		return "", err
	}

	out, err := commander.RunCommandCombined(gitCommand)
	if err != nil {
		// git command error
		return "", err
	}
	return out, nil
}

func GitConfigSet(key string, value string, gitconfig string) error {
	commandBuilder := builder.NewCommandBuilder("git")
	if gitconfig != "" {
		if !gcf.Exists(gitconfig) {
			return fmt.Errorf("gitconfig specified not found at %s", gitconfig)
		}
		commandBuilder.WithEnvironmentVariable("GIT_CONFIG", gitconfig)
	}
	gitCommand, err := commandBuilder.
		WithArguments([]string{"config", key, value}...).
		Build()

	if err != nil {
		// error creating command
		return err
	}

	_, err = commander.RunCommandCombined(gitCommand)
	if err != nil {
		// git command error
		return err
	}
	return nil
}
