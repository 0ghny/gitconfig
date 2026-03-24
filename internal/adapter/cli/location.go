package cli

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	varLocationKey      = "key"
	varShortLocationKey = "k"
)

// locationCmd returns the "location" subcommand. When called with a key as the
// first argument it prints details for that location. It also acts as the
// parent for the new and delete subcommands.
func locationCmd(factory LocationServiceFactory) *cobra.Command {
	lc := &cobra.Command{
		Use:   "location [key]",
		Short: "Manage a gitconfig location",
		Long:  `Manage a gitconfig location`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			key := args[0]
			gitconfigPath, _ := cmd.Root().PersistentFlags().GetString(varGitconfigPath)
			svc := factory(gitconfigPath)
			loc, err := svc.FindLocationByKey(key)
			if err != nil {
				log.Error(err)
				return err
			}
			if loc == nil {
				log.Debugf("location %s not found", key)
				return fmt.Errorf("location %s not found", key)
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"#", "Key", "Location", "GitConfig"})
			t.AppendRow([]interface{}{"#", loc.Key, loc.Path, loc.ConfigFile})
			t.AppendSeparator()
			t.Render()
			return nil
		},
	}

	lc.AddCommand(locationNewCmd(factory), locationDeleteCmd(factory))
	return lc
}
