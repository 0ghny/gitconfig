package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// locationDeleteCmd returns the "location delete" subcommand. It verifies the
// location exists, prompts for confirmation, then removes it from the gitconfig
// file and deletes its associated config file.
// Usage: location delete <key>
func locationDeleteCmd(factory LocationServiceFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete an existing gitconfig location",
		Long:  `Deletes an existing gitconfig location by key, removing its section from the gitconfig file and its associated config file.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			gitconfigPath, _ := cmd.Root().PersistentFlags().GetString(varGitconfigPath)
			svc := factory(gitconfigPath)

			l, err := svc.FindLocationByKey(key)
			if err != nil {
				return err
			}
			if l == nil {
				return fmt.Errorf("location '%s' not found", key)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Are you sure you want to delete location '%s'? [Y/n]: ", key)

			reader := bufio.NewReader(cmd.InOrStdin())
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read confirmation: %w", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "" && response != "y" && response != "yes" {
				fmt.Fprintln(cmd.OutOrStdout(), "Deletion cancelled.")
				return nil
			}

			if err := svc.DeleteLocation(key); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Location '%s' deleted successfully.\n", key)
			return nil
		},
	}
}
