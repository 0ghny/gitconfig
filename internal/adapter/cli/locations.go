package cli

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// locationsCmd returns the "locations" subcommand that lists all configured locations.
func locationsCmd(factory LocationServiceFactory) *cobra.Command {
	locationsListCmd := &cobra.Command{
		Use:   "locations",
		Short: "List configured locations",
		Long:  `List all configured locations with their details (it will read .gitconfig configured file)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			gitconfigPath, _ := cmd.Root().PersistentFlags().GetString(varGitconfigPath)
			svc := factory(gitconfigPath)
			locations, err := svc.GetLocations()
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"#", "Key", "Location", "GitConfig"})
			for i, l := range locations {
				t.AppendRow([]interface{}{i, l.Key, l.Path, l.ConfigFile})
				t.AppendSeparator()
			}
			t.Render()
			return nil
		},
	}
	return locationsListCmd
}
