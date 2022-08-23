package cmd

import (
	"os"

	"github.com/0ghny/gitconfigs/pkg/locations"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func locationsCmd() *cobra.Command {
	locationsListCmd := &cobra.Command{
		Use:   "locations",
		Short: "List configured locations",
		Long:  `List all configured locations with their details (it will read .gitconfig configured file)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			locationMgr := locations.NewLocationManager("", nil)
			locations, err := locationMgr.GetLocations()
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
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
