package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/0ghny/gitconfig/internal/application/locations"
	"github.com/0ghny/gitconfig/internal/domain/location"
	"github.com/0ghny/gitconfig/internal/home"
	"github.com/spf13/cobra"
)

const (
	varVerbosity          = "verbosity"
	varShortVerbosity     = "v"
	varGitconfigPath      = "git-config"
	varShortGitconfigPath = "c"
)

// LocationServiceFactory creates a location.Service for a given gitconfig path.
// Injecting it allows callers (and tests) to swap the implementation.
type LocationServiceFactory func(gitconfigPath string) location.Service

func defaultLocationServiceFactory(path string) location.Service {
	return locations.NewLocationManager(path, nil)
}

// RootCmd builds the production command tree.
func RootCmd(version string) *cobra.Command {
	return RootCmdWithDeps(version, defaultLocationServiceFactory)
}

// RootCmdWithDeps builds the command tree with injectable dependencies.
// Use this in tests to supply a custom LocationServiceFactory.
func RootCmdWithDeps(version string, factory LocationServiceFactory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "gitconfig",
		Version:      version,
		SilenceUsage: true,
		Short:        "Manage git configurations with ease",
		Long:         `Manage multiple location based git configurations easily`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := home.EnsureHome()
			return err
		},
	}

	addRootCmdFlags(rootCmd)
	rootCmd.AddCommand(locationsCmd(factory), locationCmd(factory), configCmd(factory))
	return rootCmd
}

func addRootCmdFlags(cmd *cobra.Command) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	var verbosity int
	var gitconfigPath string
	cmd.PersistentFlags().
		IntVarP(&verbosity, varVerbosity, varShortVerbosity, 0, "Verbosity level from 0 to 4")
	cmd.PersistentFlags().
		StringVarP(&gitconfigPath, varGitconfigPath, varShortGitconfigPath, filepath.Join(homeDir, ".gitconfig"), "Git configuration file")
}
