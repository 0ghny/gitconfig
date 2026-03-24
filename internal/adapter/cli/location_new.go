package cli

import (
	"fmt"

	"github.com/0ghny/gitconfig/internal/adapter/filesystem"

	"github.com/spf13/cobra"
)

const (
	varLocation      = "location"
	varShortLocation = "l"
)

// locationNewCmd returns the "location new" subcommand that creates or updates a
// gitconfig location entry in the main .gitconfig file.
// Usage: location new <key> [--location <path>]
func locationNewCmd(factory LocationServiceFactory) *cobra.Command {
	var loc string

	cmd := &cobra.Command{
		Use:   "new <key>",
		Short: "Create a new gitconfig for specified location",
		Long:  `Creates a new gitconfig file for specified location`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			gitconfigPath, _ := cmd.Root().PersistentFlags().GetString(varGitconfigPath)
			svc := factory(gitconfigPath)
			if err := svc.SaveLocation(key, loc); err != nil {
				return err
			}
			l, err := svc.FindLocationByKey(key)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Location %s saved successfully, configuration file was created at %s\n", l.Key, l.ConfigFile)
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&loc, varLocation, varShortLocation, filesystem.Getwd(), "Location for new gitconfig, default to current directory")
	return cmd
}
