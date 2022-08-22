package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/0ghny/gitconfigs/internal/home"
	"github.com/spf13/cobra"
)

const (
	varVerbosity          = "verbosity"
	varShortVerbosity     = "v"
	varGitconfigPath      = "git-config"
	varShortGitconfigPath = "c"
)

var (
	verbosity     int
	gitconfigPath string
)

func RootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gitconfigs",
		Version: version,
		Short:   "Manage git configurations with ease",
		Long:    `Manage multiple location based git configurations easily`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := home.EnsureHome()
			if err != nil {
				return err
			}
			return nil
		},
	}

	addRootCmdFlags(rootCmd)

	// Subcommands
	rootCmd.AddCommand(locationsCmd(), locationCmd())
	return rootCmd
}

func addRootCmdFlags(cmd *cobra.Command) {
	// Get user home directory to set as default for gitconfigpath
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	cmd.PersistentFlags().
		IntVarP(&verbosity, varVerbosity, varShortVerbosity, 0, "Verbosity level from 0 to 4")
	cmd.PersistentFlags().
		StringVarP(&gitconfigPath, varGitconfigPath, varShortGitconfigPath, filepath.Join(homeDir, ".gitconfig"), "Git configuration file")
}
